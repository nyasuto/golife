// 3D Game of Life WebGL Visualizer using Three.js

class LifeVisualizer {
    constructor() {
        this.scene = null;
        this.camera = null;
        this.renderer = null;
        this.controls = null;
        this.voxels = null;
        this.instancedMesh = null;
        this.ws = null;
        this.universeSize = { width: 32, height: 32, depth: 32 };
        this.fps = 0;
        this.lastTime = performance.now();
        this.frameCount = 0;

        this.init();
        this.connect();
        this.animate();
    }

    init() {
        // Scene setup
        this.scene = new THREE.Scene();
        this.scene.background = new THREE.Color(0x0a0a0a);
        this.scene.fog = new THREE.Fog(0x0a0a0a, 50, 100);

        // Camera setup
        const container = document.getElementById('canvas-container');
        this.camera = new THREE.PerspectiveCamera(
            60,
            window.innerWidth / window.innerHeight,
            0.1,
            1000
        );
        this.camera.position.set(40, 40, 40);
        this.camera.lookAt(16, 16, 16);

        // Renderer setup
        this.renderer = new THREE.WebGLRenderer({ antialias: true });
        this.renderer.setSize(window.innerWidth, window.innerHeight);
        this.renderer.setPixelRatio(window.devicePixelRatio);
        container.appendChild(this.renderer.domElement);

        // Controls setup
        this.controls = new THREE.OrbitControls(this.camera, this.renderer.domElement);
        this.controls.enableDamping = true;
        this.controls.dampingFactor = 0.05;
        this.controls.target.set(16, 16, 16);

        // Lighting
        const ambientLight = new THREE.AmbientLight(0x404040, 0.5);
        this.scene.add(ambientLight);

        const directionalLight1 = new THREE.DirectionalLight(0xffffff, 0.8);
        directionalLight1.position.set(50, 50, 50);
        this.scene.add(directionalLight1);

        const directionalLight2 = new THREE.DirectionalLight(0x4444ff, 0.3);
        directionalLight2.position.set(-50, -50, -50);
        this.scene.add(directionalLight2);

        // Grid helper
        const gridHelper = new THREE.GridHelper(32, 32, 0x333333, 0x222222);
        gridHelper.position.set(16, 0, 16);
        this.scene.add(gridHelper);

        // Bounding box
        this.createBoundingBox(32, 32, 32);

        // Create instanced mesh for voxels
        this.createInstancedMesh(1000); // Pre-allocate for up to 1000 cells

        // Window resize handler
        window.addEventListener('resize', () => this.onWindowResize(), false);
    }

    createBoundingBox(width, height, depth) {
        const edges = new THREE.EdgesGeometry(
            new THREE.BoxGeometry(width, height, depth)
        );
        const line = new THREE.LineSegments(
            edges,
            new THREE.LineBasicMaterial({ color: 0x444444 })
        );
        line.position.set(width / 2, height / 2, depth / 2);
        this.scene.add(line);
    }

    createInstancedMesh(maxInstances) {
        const geometry = new THREE.BoxGeometry(0.9, 0.9, 0.9);
        const material = new THREE.MeshPhongMaterial({
            color: 0x00ff00,
            emissive: 0x002200,
            specular: 0x111111,
            shininess: 30
        });

        this.instancedMesh = new THREE.InstancedMesh(geometry, material, maxInstances);
        this.instancedMesh.instanceMatrix.setUsage(THREE.DynamicDrawUsage);
        this.scene.add(this.instancedMesh);
    }

    connect() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws`;

        console.log('Connecting to WebSocket:', wsUrl);
        this.ws = new WebSocket(wsUrl);

        this.ws.onopen = () => {
            console.log('WebSocket connected');
            this.updateConnectionStatus(true);
        };

        this.ws.onclose = () => {
            console.log('WebSocket disconnected');
            this.updateConnectionStatus(false);
            // Attempt reconnection after 2 seconds
            setTimeout(() => this.connect(), 2000);
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        this.ws.onmessage = (event) => {
            const state = JSON.parse(event.data);
            this.updateVisualization(state);
        };
    }

    updateVisualization(state) {
        // Update info panel
        document.getElementById('generation').textContent = state.generation;
        document.getElementById('population').textContent = state.population;
        document.getElementById('universe-size').textContent =
            `${state.width}×${state.height}×${state.depth}`;

        // Update universe size if changed
        if (state.width !== this.universeSize.width ||
            state.height !== this.universeSize.height ||
            state.depth !== this.universeSize.depth) {
            this.universeSize = {
                width: state.width,
                height: state.height,
                depth: state.depth
            };
            this.controls.target.set(state.width / 2, state.height / 2, state.depth / 2);
        }

        // Update voxels using instanced rendering
        this.updateVoxels(state.cells);
    }

    updateVoxels(cells) {
        const matrix = new THREE.Matrix4();
        const color = new THREE.Color();

        // Set instance count
        this.instancedMesh.count = cells.length;

        // Update each instance
        cells.forEach((cell, index) => {
            // Position matrix
            matrix.setPosition(cell.x + 0.5, cell.y + 0.5, cell.z + 0.5);
            this.instancedMesh.setMatrixAt(index, matrix);

            // Color based on position (optional: could be based on age)
            const hue = (cell.z / this.universeSize.depth) * 0.3 + 0.3; // Green to cyan gradient
            color.setHSL(hue, 1.0, 0.5);
            this.instancedMesh.setColorAt(index, color);
        });

        // Update the instance matrix
        this.instancedMesh.instanceMatrix.needsUpdate = true;
        if (this.instancedMesh.instanceColor) {
            this.instancedMesh.instanceColor.needsUpdate = true;
        }
    }

    updateConnectionStatus(connected) {
        const statusEl = document.getElementById('connection-status');
        if (connected) {
            statusEl.textContent = '● Connected';
            statusEl.className = 'connected';
        } else {
            statusEl.textContent = '● Disconnected';
            statusEl.className = 'disconnected';
        }
    }

    onWindowResize() {
        this.camera.aspect = window.innerWidth / window.innerHeight;
        this.camera.updateProjectionMatrix();
        this.renderer.setSize(window.innerWidth, window.innerHeight);
    }

    updateFPS() {
        this.frameCount++;
        const currentTime = performance.now();
        const elapsed = currentTime - this.lastTime;

        if (elapsed >= 1000) {
            this.fps = Math.round((this.frameCount * 1000) / elapsed);
            document.getElementById('fps').textContent = this.fps;
            this.frameCount = 0;
            this.lastTime = currentTime;
        }
    }

    animate() {
        requestAnimationFrame(() => this.animate());

        this.controls.update();
        this.renderer.render(this.scene, this.camera);
        this.updateFPS();
    }
}

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    new LifeVisualizer();
});
