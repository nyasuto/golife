/**
 * 3D Game of Life WASM Viewer
 * Combines Three.js visualization with Go WASM backend
 */

class Life3DViewer {
    constructor() {
        this.scene = null;
        this.camera = null;
        this.renderer = null;
        this.controls = null;
        this.instancedMesh = null;
        this.wasmAPI = null;
        this.universeSize = { width: 10, height: 10, depth: 10 };
        this.animationController = null;
        this.isAnimating = false;

        // Performance tracking
        this.fps = 0;
        this.lastTime = performance.now();
        this.frameCount = 0;

        this.init();
    }

    async init() {
        // Initialize WASM
        this.updateStatus('Loading WASM...', 'loading');

        this.wasmAPI = new GoLifeWASM();

        try {
            await this.wasmAPI.load('life3d.wasm');
            this.updateStatus('✅ WASM Ready', 'ready');
            console.log('WASM API loaded successfully');
        } catch (err) {
            this.updateStatus('❌ WASM Load Failed', 'error');
            console.error('Failed to load WASM:', err);
            return;
        }

        // Initialize Three.js
        this.initThreeJS();

        // Initialize with default 10x10x10 universe
        this.initUniverse();

        // Start render loop
        this.animate();
    }

    initThreeJS() {
        // Scene setup
        this.scene = new THREE.Scene();
        this.scene.background = new THREE.Color(0x0a0a0a);
        this.scene.fog = new THREE.Fog(0x0a0a0a, 20, 60);

        // Camera setup
        const container = document.getElementById('canvas-container');
        this.camera = new THREE.PerspectiveCamera(
            60,
            window.innerWidth / window.innerHeight,
            0.1,
            1000
        );
        this.camera.position.set(15, 15, 15);
        this.camera.lookAt(5, 5, 5);

        // Renderer setup
        this.renderer = new THREE.WebGLRenderer({ antialias: true });
        this.renderer.setSize(window.innerWidth, window.innerHeight);
        this.renderer.setPixelRatio(window.devicePixelRatio);
        container.appendChild(this.renderer.domElement);

        // Controls setup
        this.controls = new THREE.OrbitControls(this.camera, this.renderer.domElement);
        this.controls.enableDamping = true;
        this.controls.dampingFactor = 0.05;
        this.controls.target.set(5, 5, 5);

        // Lighting
        const ambientLight = new THREE.AmbientLight(0x404040, 0.6);
        this.scene.add(ambientLight);

        const directionalLight1 = new THREE.DirectionalLight(0xffffff, 0.8);
        directionalLight1.position.set(30, 30, 30);
        this.scene.add(directionalLight1);

        const directionalLight2 = new THREE.DirectionalLight(0x4444ff, 0.4);
        directionalLight2.position.set(-30, -30, -30);
        this.scene.add(directionalLight2);

        // Grid helper
        const gridHelper = new THREE.GridHelper(10, 10, 0x333333, 0x222222);
        gridHelper.position.set(5, 0, 5);
        this.scene.add(gridHelper);

        // Create bounding box
        this.createBoundingBox(10, 10, 10);

        // Create instanced mesh for voxels
        this.createInstancedMesh(200); // Pre-allocate for up to 200 cells

        // Window resize handler
        window.addEventListener('resize', () => this.onWindowResize(), false);
    }

    createBoundingBox(width, height, depth) {
        // Remove old bounding box if exists
        const oldBox = this.scene.getObjectByName('boundingBox');
        if (oldBox) {
            this.scene.remove(oldBox);
        }

        const edges = new THREE.EdgesGeometry(
            new THREE.BoxGeometry(width, height, depth)
        );
        const line = new THREE.LineSegments(
            edges,
            new THREE.LineBasicMaterial({ color: 0x00d4ff, linewidth: 2 })
        );
        line.position.set(width / 2, height / 2, depth / 2);
        line.name = 'boundingBox';
        this.scene.add(line);

        // Remove old grid
        const oldGrid = this.scene.getObjectByName('gridHelper');
        if (oldGrid) {
            this.scene.remove(oldGrid);
        }

        // Update grid
        const gridHelper = new THREE.GridHelper(width, width, 0x333333, 0x222222);
        gridHelper.position.set(width / 2, 0, depth / 2);
        gridHelper.name = 'gridHelper';
        this.scene.add(gridHelper);
    }

    createInstancedMesh(maxInstances) {
        // Remove old mesh if exists
        if (this.instancedMesh) {
            this.scene.remove(this.instancedMesh);
        }

        const geometry = new THREE.BoxGeometry(0.8, 0.8, 0.8);
        const material = new THREE.MeshPhongMaterial({
            color: 0x00ff88,
            emissive: 0x002200,
            specular: 0x111111,
            shininess: 50,
            transparent: true,
            opacity: 0.9
        });

        this.instancedMesh = new THREE.InstancedMesh(geometry, material, maxInstances);
        this.instancedMesh.instanceMatrix.setUsage(THREE.DynamicDrawUsage);
        this.scene.add(this.instancedMesh);
    }

    initUniverse() {
        const sizeSelect = document.getElementById('universeSize');
        const size = parseInt(sizeSelect.value);

        console.log(`Initializing universe: ${size}×${size}×${size}`);
        const result = this.wasmAPI.initUniverse(size, size, size);

        if (result.error) {
            console.error('Failed to initialize universe:', result.error);
            this.updateStatus('❌ Init Failed', 'error');
            return;
        }

        this.universeSize = { width: size, height: size, depth: size };

        // Update camera and controls
        const center = size / 2;
        this.controls.target.set(center, center, center);
        this.camera.position.set(size * 1.5, size * 1.5, size * 1.5);
        this.camera.lookAt(center, center, center);

        // Update bounding box
        this.createBoundingBox(size, size, size);

        // Recreate instanced mesh with appropriate capacity
        const maxCells = Math.min(size * size * size / 2, 1000); // Estimate max cells
        this.createInstancedMesh(maxCells);

        this.updateInfo();
        this.updateVisualization();

        console.log('Universe initialized successfully');
    }

    loadPattern() {
        const patternSelect = document.getElementById('pattern');
        const pattern = patternSelect.value;
        const center = Math.floor(this.universeSize.width / 2);

        console.log(`Loading pattern: ${pattern} at (${center}, ${center}, ${center})`);
        const result = this.wasmAPI.loadPattern(pattern, center - 2, center - 2, center - 2);

        if (result.error) {
            console.error('Failed to load pattern:', result.error);
            return;
        }

        this.updateInfo();
        this.updateVisualization();

        console.log(`Pattern '${pattern}' loaded successfully`);
    }

    step() {
        const result = this.wasmAPI.step();

        if (result.error) {
            console.error('Step error:', result.error);
            return;
        }

        this.updateInfo();
        this.updateVisualization();
    }

    toggleAnimation() {
        const btn = document.getElementById('animateBtn');

        if (this.isAnimating) {
            // Stop animation
            if (this.animationController) {
                this.animationController.stop();
                this.animationController = null;
            }
            this.isAnimating = false;
            btn.textContent = 'Start';
            btn.style.background = 'linear-gradient(135deg, #00d4ff, #0088ff)';
            console.log('Animation stopped');
        } else {
            // Start animation
            this.animationController = this.wasmAPI.startAnimation((state, generation) => {
                this.updateInfo();
                this.updateVisualization();
            }, 10); // 10 FPS

            this.isAnimating = true;
            btn.textContent = 'Stop';
            btn.style.background = 'linear-gradient(135deg, #ff6b6b, #ee5a6f)';
            console.log('Animation started at 10 FPS');
        }
    }

    clearUniverse() {
        if (this.isAnimating) {
            this.toggleAnimation();
        }

        const result = this.wasmAPI.clearUniverse();

        if (result.error) {
            console.error('Clear error:', result.error);
            return;
        }

        this.updateInfo();
        this.updateVisualization();

        console.log('Universe cleared');
    }

    updateVisualization() {
        const state = this.wasmAPI.getLivingCells();

        if (state.error) {
            console.error('Failed to get living cells:', state.error);
            return;
        }

        const cells = state.cells || [];
        const matrix = new THREE.Matrix4();

        // Update instanced mesh
        for (let i = 0; i < cells.length; i++) {
            const cell = cells[i];
            matrix.setPosition(cell.x + 0.5, cell.y + 0.5, cell.z + 0.5);
            this.instancedMesh.setMatrixAt(i, matrix);
        }

        // Hide unused instances
        matrix.setPosition(1000, 1000, 1000); // Move far away
        for (let i = cells.length; i < this.instancedMesh.count; i++) {
            this.instancedMesh.setMatrixAt(i, matrix);
        }

        this.instancedMesh.instanceMatrix.needsUpdate = true;
        this.instancedMesh.count = Math.max(cells.length, 0);
    }

    updateInfo() {
        const info = this.wasmAPI.getUniverseInfo();

        if (info.error) {
            console.error('Failed to get universe info:', info.error);
            return;
        }

        document.getElementById('generation').textContent = info.generation || 0;
        document.getElementById('population').textContent = info.population || 0;
        document.getElementById('universe-size').textContent =
            `${info.width}×${info.height}×${info.depth}`;
    }

    updateStatus(text, status) {
        const statusDiv = document.getElementById('status');
        statusDiv.textContent = text;
        statusDiv.className = status;
    }

    onWindowResize() {
        this.camera.aspect = window.innerWidth / window.innerHeight;
        this.camera.updateProjectionMatrix();
        this.renderer.setSize(window.innerWidth, window.innerHeight);
    }

    animate() {
        requestAnimationFrame(() => this.animate());

        // Update controls
        this.controls.update();

        // Calculate FPS
        const currentTime = performance.now();
        this.frameCount++;

        if (currentTime - this.lastTime >= 1000) {
            this.fps = Math.round(this.frameCount / ((currentTime - this.lastTime) / 1000));
            document.getElementById('fps').textContent = this.fps;
            this.frameCount = 0;
            this.lastTime = currentTime;
        }

        // Render scene
        this.renderer.render(this.scene, this.camera);
    }
}

// Global app instance
let app;

// Initialize when page loads
window.addEventListener('load', () => {
    app = new Life3DViewer();
});
