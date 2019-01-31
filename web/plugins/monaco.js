window.MonacoEnvironment = {
    getWorkerUrl: function (moduleId, label) {
        if (label === 'json') {
            return '/monaco/json.worker.bundle.js';
        }
        if (label === 'typescript' || label === 'javascript') {
            return '/monaco/ts.worker.bundle.js';
        }
        return '/monaco/editor.worker.bundle.js';
    }
}
