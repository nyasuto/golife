/**
 * GoLife WASM Communication Layer
 *
 * Provides a JavaScript API for interacting with the Go WASM Game of Life simulation
 */

class GoLifeWASM {
    constructor() {
        this.wasmReady = false;
        this.onReadyCallbacks = [];
    }

    /**
     * Load the WASM module
     * @param {string} wasmPath - Path to the .wasm file (default: 'life3d.wasm')
     * @returns {Promise<void>}
     */
    async load(wasmPath = 'life3d.wasm') {
        try {
            console.log('🔄 Loading Go WASM module...');

            // Create Go runtime
            const go = new Go();

            // Load and instantiate WASM
            const result = await WebAssembly.instantiateStreaming(
                fetch(wasmPath),
                go.importObject
            );

            console.log('✅ WASM module loaded');

            // Run Go program (this will register callbacks)
            go.run(result.instance);

            // Wait a bit for callbacks to be registered
            await new Promise(resolve => setTimeout(resolve, 100));

            // Verify callbacks are available
            if (typeof window.goInitUniverse !== 'function') {
                throw new Error('Go functions not registered properly');
            }

            this.wasmReady = true;
            console.log('✅ Go WASM ready');

            // Call ready callbacks
            this.onReadyCallbacks.forEach(callback => callback());
            this.onReadyCallbacks = [];

        } catch (err) {
            console.error('❌ WASM load failed:', err);
            throw err;
        }
    }

    /**
     * Register a callback to be called when WASM is ready
     * @param {Function} callback
     */
    onReady(callback) {
        if (this.wasmReady) {
            callback();
        } else {
            this.onReadyCallbacks.push(callback);
        }
    }

    /**
     * Check if WASM is ready
     * @returns {boolean}
     */
    isReady() {
        return this.wasmReady;
    }

    /**
     * Initialize a 3D universe
     * @param {number} width - Universe width
     * @param {number} height - Universe height
     * @param {number} depth - Universe depth
     * @returns {Object} Result object with success/error
     */
    initUniverse(width, height, depth) {
        if (!this.wasmReady) {
            return { error: 'WASM not ready' };
        }
        return window.goInitUniverse(width, height, depth);
    }

    /**
     * Load a pattern into the universe
     * @param {string} patternName - Pattern name ('glider', 'block', 'tube')
     * @param {number} x - X position
     * @param {number} y - Y position
     * @param {number} z - Z position
     * @returns {Object} Result object with success/error
     */
    loadPattern(patternName, x, y, z) {
        if (!this.wasmReady) {
            return { error: 'WASM not ready' };
        }
        return window.goLoadPattern(patternName, x, y, z);
    }

    /**
     * Advance simulation by one generation
     * @returns {Object} Result object with success/error and generation number
     */
    step() {
        if (!this.wasmReady) {
            return { error: 'WASM not ready' };
        }
        return window.goStep();
    }

    /**
     * Get all living cells
     * @returns {Object|string} Universe state as JSON string or error object
     */
    getLivingCells() {
        if (!this.wasmReady) {
            return { error: 'WASM not ready' };
        }
        const result = window.goGetLivingCells();

        // Parse JSON if it's a string
        if (typeof result === 'string') {
            return JSON.parse(result);
        }
        return result;
    }

    /**
     * Get universe information
     * @returns {Object} Universe info (width, height, depth, generation, population)
     */
    getUniverseInfo() {
        if (!this.wasmReady) {
            return { error: 'WASM not ready' };
        }
        return window.goGetUniverseInfo();
    }

    /**
     * Clear the universe (set all cells to dead)
     * @returns {Object} Result object with success/error
     */
    clearUniverse() {
        if (!this.wasmReady) {
            return { error: 'WASM not ready' };
        }
        return window.goClearUniverse();
    }

    /**
     * Set a specific cell's state
     * @param {number} x - X coordinate
     * @param {number} y - Y coordinate
     * @param {number} z - Z coordinate
     * @param {boolean} alive - Cell state
     * @returns {Object} Result object with success/error
     */
    setCell(x, y, z, alive) {
        if (!this.wasmReady) {
            return { error: 'WASM not ready' };
        }
        return window.goSetCell(x, y, z, alive);
    }

    /**
     * Run simulation for N generations
     * @param {number} generations - Number of generations to run
     * @param {Function} onStep - Callback called after each step (optional)
     * @returns {Promise<void>}
     */
    async run(generations, onStep = null) {
        if (!this.wasmReady) {
            throw new Error('WASM not ready');
        }

        for (let i = 0; i < generations; i++) {
            const result = this.step();
            if (result.error) {
                throw new Error(result.error);
            }

            if (onStep) {
                const state = this.getLivingCells();
                await onStep(state, i + 1);
            }
        }
    }

    /**
     * Start continuous simulation with animation
     * @param {Function} onFrame - Callback called each frame with (state, generation)
     * @param {number} fps - Frames per second (default: 10)
     * @returns {Object} Controller object with stop() method
     */
    startAnimation(onFrame, fps = 10) {
        if (!this.wasmReady) {
            throw new Error('WASM not ready');
        }

        const interval = 1000 / fps;
        let running = true;

        const animate = () => {
            if (!running) return;

            const stepResult = this.step();
            if (stepResult.error) {
                console.error('Step error:', stepResult.error);
                running = false;
                return;
            }

            const state = this.getLivingCells();
            if (state.error) {
                console.error('GetLivingCells error:', state.error);
                running = false;
                return;
            }

            onFrame(state, stepResult.generation);

            setTimeout(animate, interval);
        };

        setTimeout(animate, interval);

        return {
            stop: () => {
                running = false;
            }
        };
    }
}

// Create global instance
window.GoLifeWASM = GoLifeWASM;
