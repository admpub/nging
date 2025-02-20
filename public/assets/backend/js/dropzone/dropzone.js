// modules are defined as an array
// [ module function, map of requires ]
//
// map of requires is short require name -> numeric require
//
// anything defined in a previous bundle is accessed via the
// orig method which is the require for previous bundles

(function (modules, entry, mainEntry, parcelRequireName, globalName) {
  /* eslint-disable no-undef */
  var globalObject =
    typeof globalThis !== 'undefined'
      ? globalThis
      : typeof self !== 'undefined'
      ? self
      : typeof window !== 'undefined'
      ? window
      : typeof global !== 'undefined'
      ? global
      : {};
  /* eslint-enable no-undef */

  // Save the require from previous bundle to this closure if any
  var previousRequire =
    typeof globalObject[parcelRequireName] === 'function' &&
    globalObject[parcelRequireName];

  var cache = previousRequire.cache || {};
  // Do not use `require` to prevent Webpack from trying to bundle this call
  var nodeRequire =
    typeof module !== 'undefined' &&
    typeof module.require === 'function' &&
    module.require.bind(module);

  function newRequire(name, jumped) {
    if (!cache[name]) {
      if (!modules[name]) {
        // if we cannot find the module within our internal map or
        // cache jump to the current global require ie. the last bundle
        // that was added to the page.
        var currentRequire =
          typeof globalObject[parcelRequireName] === 'function' &&
          globalObject[parcelRequireName];
        if (!jumped && currentRequire) {
          return currentRequire(name, true);
        }

        // If there are other bundles on this page the require from the
        // previous one is saved to 'previousRequire'. Repeat this as
        // many times as there are bundles until the module is found or
        // we exhaust the require chain.
        if (previousRequire) {
          return previousRequire(name, true);
        }

        // Try the node require function if it exists.
        if (nodeRequire && typeof name === 'string') {
          return nodeRequire(name);
        }

        var err = new Error("Cannot find module '" + name + "'");
        err.code = 'MODULE_NOT_FOUND';
        throw err;
      }

      localRequire.resolve = resolve;
      localRequire.cache = {};

      var module = (cache[name] = new newRequire.Module(name));

      modules[name][0].call(
        module.exports,
        localRequire,
        module,
        module.exports,
        this
      );
    }

    return cache[name].exports;

    function localRequire(x) {
      var res = localRequire.resolve(x);
      return res === false ? {} : newRequire(res);
    }

    function resolve(x) {
      var id = modules[name][1][x];
      return id != null ? id : x;
    }
  }

  function Module(moduleName) {
    this.id = moduleName;
    this.bundle = newRequire;
    this.exports = {};
  }

  newRequire.isParcelRequire = true;
  newRequire.Module = Module;
  newRequire.modules = modules;
  newRequire.cache = cache;
  newRequire.parent = previousRequire;
  newRequire.register = function (id, exports) {
    modules[id] = [
      function (require, module) {
        module.exports = exports;
      },
      {},
    ];
  };

  Object.defineProperty(newRequire, 'root', {
    get: function () {
      return globalObject[parcelRequireName];
    },
  });

  globalObject[parcelRequireName] = newRequire;

  for (var i = 0; i < entry.length; i++) {
    newRequire(entry[i]);
  }

  if (mainEntry) {
    // Expose entry point to Node, AMD or browser globals
    // Based on https://github.com/ForbesLindesay/umd/blob/master/template.js
    var mainExports = newRequire(mainEntry);

    // CommonJS
    if (typeof exports === 'object' && typeof module !== 'undefined') {
      module.exports = mainExports;

      // RequireJS
    } else if (typeof define === 'function' && define.amd) {
      define(function () {
        return mainExports;
      });

      // <script>
    } else if (globalName) {
      this[globalName] = mainExports;
    }
  }
})({"hq6rc":[function(require,module,exports) {
var _asyncToGenerator = require("@swc/helpers/_/_async_to_generator");
var _toConsumableArray = require("@swc/helpers/_/_to_consumable_array");
var _tsGenerator = require("@swc/helpers/_/_ts_generator");
var global = arguments[3];
var HMR_HOST = null;
var HMR_PORT = 1234;
var HMR_SECURE = false;
var HMR_ENV_HASH = "891821638b727f7d";
var HMR_USE_SSE = false;
module.bundle.HMR_BUNDLE_ID = "887532d8794fcd48";
"use strict";
/* global HMR_HOST, HMR_PORT, HMR_ENV_HASH, HMR_SECURE, HMR_USE_SSE, chrome, browser, __parcel__import__, __parcel__importScripts__, ServiceWorkerGlobalScope */ /*::
import type {
  HMRAsset,
  HMRMessage,
} from '@parcel/reporter-dev-server/src/HMRServer.js';
interface ParcelRequire {
  (string): mixed;
  cache: {|[string]: ParcelModule|};
  hotData: {|[string]: mixed|};
  Module: any;
  parent: ?ParcelRequire;
  isParcelRequire: true;
  modules: {|[string]: [Function, {|[string]: string|}]|};
  HMR_BUNDLE_ID: string;
  root: ParcelRequire;
}
interface ParcelModule {
  hot: {|
    data: mixed,
    accept(cb: (Function) => void): void,
    dispose(cb: (mixed) => void): void,
    // accept(deps: Array<string> | string, cb: (Function) => void): void,
    // decline(): void,
    _acceptCallbacks: Array<(Function) => void>,
    _disposeCallbacks: Array<(mixed) => void>,
  |};
}
interface ExtensionContext {
  runtime: {|
    reload(): void,
    getURL(url: string): string;
    getManifest(): {manifest_version: number, ...};
  |};
}
declare var module: {bundle: ParcelRequire, ...};
declare var HMR_HOST: string;
declare var HMR_PORT: string;
declare var HMR_ENV_HASH: string;
declare var HMR_SECURE: boolean;
declare var HMR_USE_SSE: boolean;
declare var chrome: ExtensionContext;
declare var browser: ExtensionContext;
declare var __parcel__import__: (string) => Promise<void>;
declare var __parcel__importScripts__: (string) => Promise<void>;
declare var globalThis: typeof self;
declare var ServiceWorkerGlobalScope: Object;
*/ var OVERLAY_ID = "__parcel__error__overlay__";
var OldModule = module.bundle.Module;
function Module(moduleName) {
    OldModule.call(this, moduleName);
    this.hot = {
        data: module.bundle.hotData[moduleName],
        _acceptCallbacks: [],
        _disposeCallbacks: [],
        accept: function accept(fn) {
            this._acceptCallbacks.push(fn || function() {});
        },
        dispose: function dispose(fn) {
            this._disposeCallbacks.push(fn);
        }
    };
    module.bundle.hotData[moduleName] = undefined;
}
module.bundle.Module = Module;
module.bundle.hotData = {};
var checkedAssets /*: {|[string]: boolean|} */ , assetsToDispose /*: Array<[ParcelRequire, string]> */ , assetsToAccept /*: Array<[ParcelRequire, string]> */ ;
function getHostname() {
    return HMR_HOST || (location.protocol.indexOf("http") === 0 ? location.hostname : "localhost");
}
function getPort() {
    return HMR_PORT || location.port;
}
// eslint-disable-next-line no-redeclare
var parent = module.bundle.parent;
if ((!parent || !parent.isParcelRequire) && typeof WebSocket !== "undefined") {
    var hostname = getHostname();
    var port = getPort();
    var protocol = HMR_SECURE || location.protocol == "https:" && ![
        "localhost",
        "127.0.0.1",
        "0.0.0.0"
    ].includes(hostname) ? "wss" : "ws";
    var ws;
    if (HMR_USE_SSE) ws = new EventSource("/__parcel_hmr");
    else try {
        ws = new WebSocket(protocol + "://" + hostname + (port ? ":" + port : "") + "/");
    } catch (err) {
        if (err.message) console.error(err.message);
        ws = {};
    }
    // Web extension context
    var extCtx = typeof browser === "undefined" ? typeof chrome === "undefined" ? null : chrome : browser;
    // Safari doesn't support sourceURL in error stacks.
    // eval may also be disabled via CSP, so do a quick check.
    var supportsSourceURL = false;
    try {
        (0, eval)('throw new Error("test"); //# sourceURL=test.js');
    } catch (err) {
        supportsSourceURL = err.stack.includes("test.js");
    }
    // $FlowFixMe
    ws.onmessage = function() {
        var _ref = (0, _asyncToGenerator._)(function(event /*: {data: string, ...} */ ) {
            var data /*: HMRMessage */ , assets, handled, processedAssets, i, id, i1, id1, _iteratorNormalCompletion, _didIteratorError, _iteratorError, _iterator, _step, ansiDiagnostic, stack, overlay;
            return (0, _tsGenerator._)(this, function(_state) {
                switch(_state.label){
                    case 0:
                        checkedAssets = {} /*: {|[string]: boolean|} */ ;
                        assetsToAccept = [];
                        assetsToDispose = [];
                        data = JSON.parse(event.data);
                        if (!(data.type === "update")) return [
                            3,
                            3
                        ];
                        // Remove error overlay if there is one
                        if (typeof document !== "undefined") removeErrorOverlay();
                        assets = data.assets.filter(function(asset) {
                            return asset.envHash === HMR_ENV_HASH;
                        });
                        // Handle HMR Update
                        handled = assets.every(function(asset) {
                            return asset.type === "css" || asset.type === "js" && hmrAcceptCheck(module.bundle.root, asset.id, asset.depsByBundle);
                        });
                        if (!handled) return [
                            3,
                            2
                        ];
                        console.clear();
                        // Dispatch custom event so other runtimes (e.g React Refresh) are aware.
                        if (typeof window !== "undefined" && typeof CustomEvent !== "undefined") window.dispatchEvent(new CustomEvent("parcelhmraccept"));
                        return [
                            4,
                            hmrApplyUpdates(assets)
                        ];
                    case 1:
                        _state.sent();
                        // Dispose all old assets.
                        processedAssets = {} /*: {|[string]: boolean|} */ ;
                        for(i = 0; i < assetsToDispose.length; i++){
                            id = assetsToDispose[i][1];
                            if (!processedAssets[id]) {
                                hmrDispose(assetsToDispose[i][0], id);
                                processedAssets[id] = true;
                            }
                        }
                        // Run accept callbacks. This will also re-execute other disposed assets in topological order.
                        processedAssets = {};
                        for(i1 = 0; i1 < assetsToAccept.length; i1++){
                            id1 = assetsToAccept[i1][1];
                            if (!processedAssets[id1]) {
                                hmrAccept(assetsToAccept[i1][0], id1);
                                processedAssets[id1] = true;
                            }
                        }
                        return [
                            3,
                            3
                        ];
                    case 2:
                        fullReload();
                        _state.label = 3;
                    case 3:
                        if (data.type === "error") {
                            _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                            try {
                                // Log parcel errors to console
                                for(_iterator = data.diagnostics.ansi[Symbol.iterator](); !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                                    ansiDiagnostic = _step.value;
                                    stack = ansiDiagnostic.codeframe ? ansiDiagnostic.codeframe : ansiDiagnostic.stack;
                                    console.error("\uD83D\uDEA8 [parcel]: " + ansiDiagnostic.message + "\n" + stack + "\n\n" + ansiDiagnostic.hints.join("\n"));
                                }
                            } catch (err) {
                                _didIteratorError = true;
                                _iteratorError = err;
                            } finally{
                                try {
                                    if (!_iteratorNormalCompletion && _iterator.return != null) {
                                        _iterator.return();
                                    }
                                } finally{
                                    if (_didIteratorError) {
                                        throw _iteratorError;
                                    }
                                }
                            }
                            if (typeof document !== "undefined") {
                                // Render the fancy html overlay
                                removeErrorOverlay();
                                overlay = createErrorOverlay(data.diagnostics.html);
                                // $FlowFixMe
                                document.body.appendChild(overlay);
                            }
                        }
                        return [
                            2
                        ];
                }
            });
        });
        return function(event) {
            return _ref.apply(this, arguments);
        };
    }();
    if (ws instanceof WebSocket) {
        ws.onerror = function(e) {
            if (e.message) console.error(e.message);
        };
        ws.onclose = function() {
            console.warn("[parcel] \uD83D\uDEA8 Connection to the HMR server was lost");
        };
    }
}
function removeErrorOverlay() {
    var overlay = document.getElementById(OVERLAY_ID);
    if (overlay) {
        overlay.remove();
        console.log("[parcel] \u2728 Error resolved");
    }
}
function createErrorOverlay(diagnostics) {
    var overlay = document.createElement("div");
    overlay.id = OVERLAY_ID;
    var errorHTML = '<div style="background: black; opacity: 0.85; font-size: 16px; color: white; position: fixed; height: 100%; width: 100%; top: 0px; left: 0px; padding: 30px; font-family: Menlo, Consolas, monospace; z-index: 9999;">';
    var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
    try {
        for(var _iterator = diagnostics[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
            var diagnostic = _step.value;
            var stack = diagnostic.frames.length ? diagnostic.frames.reduce(function(p, frame) {
                return "".concat(p, '\n<a href="/__parcel_launch_editor?file=').concat(encodeURIComponent(frame.location), '" style="text-decoration: underline; color: #888" onclick="fetch(this.href); return false">').concat(frame.location, "</a>\n").concat(frame.code);
            }, "") : diagnostic.stack;
            errorHTML += '\n      <div>\n        <div style="font-size: 18px; font-weight: bold; margin-top: 20px;">\n          \uD83D\uDEA8 '.concat(diagnostic.message, "\n        </div>\n        <pre>").concat(stack, "</pre>\n        <div>\n          ").concat(diagnostic.hints.map(function(hint) {
                return "<div>\uD83D\uDCA1 " + hint + "</div>";
            }).join(""), "\n        </div>\n        ").concat(diagnostic.documentation ? '<div>\uD83D\uDCDD <a style="color: violet" href="'.concat(diagnostic.documentation, '" target="_blank">Learn more</a></div>') : "", "\n      </div>\n    ");
        }
    } catch (err) {
        _didIteratorError = true;
        _iteratorError = err;
    } finally{
        try {
            if (!_iteratorNormalCompletion && _iterator.return != null) {
                _iterator.return();
            }
        } finally{
            if (_didIteratorError) {
                throw _iteratorError;
            }
        }
    }
    errorHTML += "</div>";
    overlay.innerHTML = errorHTML;
    return overlay;
}
function fullReload() {
    if ("reload" in location) location.reload();
    else if (extCtx && extCtx.runtime && extCtx.runtime.reload) extCtx.runtime.reload();
}
function getParents(bundle, id) /*: Array<[ParcelRequire, string]> */ {
    var modules = bundle.modules;
    if (!modules) return [];
    var parents = [];
    var k, d, dep;
    for(k in modules)for(d in modules[k][1]){
        dep = modules[k][1][d];
        if (dep === id || Array.isArray(dep) && dep[dep.length - 1] === id) parents.push([
            bundle,
            k
        ]);
    }
    if (bundle.parent) parents = parents.concat(getParents(bundle.parent, id));
    return parents;
}
function updateLink(link) {
    var href = link.getAttribute("href");
    if (!href) return;
    var newLink = link.cloneNode();
    newLink.onload = function() {
        if (link.parentNode !== null) // $FlowFixMe
        link.parentNode.removeChild(link);
    };
    newLink.setAttribute("href", // $FlowFixMe
    href.split("?")[0] + "?" + Date.now());
    // $FlowFixMe
    link.parentNode.insertBefore(newLink, link.nextSibling);
}
var cssTimeout = null;
function reloadCSS() {
    if (cssTimeout) return;
    cssTimeout = setTimeout(function() {
        var links = document.querySelectorAll('link[rel="stylesheet"]');
        for(var i = 0; i < links.length; i++){
            // $FlowFixMe[incompatible-type]
            var href /*: string */  = links[i].getAttribute("href");
            var hostname = getHostname();
            var servedFromHMRServer = hostname === "localhost" ? new RegExp("^(https?:\\/\\/(0.0.0.0|127.0.0.1)|localhost):" + getPort()).test(href) : href.indexOf(hostname + ":" + getPort());
            var absolute = /^https?:\/\//i.test(href) && href.indexOf(location.origin) !== 0 && !servedFromHMRServer;
            if (!absolute) updateLink(links[i]);
        }
        cssTimeout = null;
    }, 50);
}
function hmrDownload(asset) {
    if (asset.type === "js") {
        if (typeof document !== "undefined") {
            var script = document.createElement("script");
            script.src = asset.url + "?t=" + Date.now();
            if (asset.outputFormat === "esmodule") script.type = "module";
            return new Promise(function(resolve, reject) {
                var _document$head;
                script.onload = function() {
                    return resolve(script);
                };
                script.onerror = reject;
                (_document$head = document.head) === null || _document$head === void 0 || _document$head.appendChild(script);
            });
        } else if (typeof importScripts === "function") {
            // Worker scripts
            if (asset.outputFormat === "esmodule") return import(asset.url + "?t=" + Date.now());
            else return new Promise(function(resolve, reject) {
                try {
                    importScripts(asset.url + "?t=" + Date.now());
                    resolve();
                } catch (err) {
                    reject(err);
                }
            });
        }
    }
}
function hmrApplyUpdates(assets) {
    return _hmrApplyUpdates.apply(this, arguments);
}
function _hmrApplyUpdates() {
    _hmrApplyUpdates = (0, _asyncToGenerator._)(function(assets) {
        var scriptsToRemove, promises;
        return (0, _tsGenerator._)(this, function(_state) {
            switch(_state.label){
                case 0:
                    global.parcelHotUpdate = Object.create(null);
                    _state.label = 1;
                case 1:
                    _state.trys.push([
                        1,
                        ,
                        4,
                        5
                    ]);
                    if (!!supportsSourceURL) return [
                        3,
                        3
                    ];
                    promises = assets.map(function(asset) {
                        var _hmrDownload;
                        return (_hmrDownload = hmrDownload(asset)) === null || _hmrDownload === void 0 ? void 0 : _hmrDownload.catch(function(err) {
                            // Web extension fix
                            if (extCtx && extCtx.runtime && extCtx.runtime.getManifest().manifest_version == 3 && typeof ServiceWorkerGlobalScope != "undefined" && global instanceof ServiceWorkerGlobalScope) {
                                extCtx.runtime.reload();
                                return;
                            }
                            throw err;
                        });
                    });
                    return [
                        4,
                        Promise.all(promises)
                    ];
                case 2:
                    scriptsToRemove = _state.sent();
                    _state.label = 3;
                case 3:
                    assets.forEach(function(asset) {
                        hmrApply(module.bundle.root, asset);
                    });
                    return [
                        3,
                        5
                    ];
                case 4:
                    delete global.parcelHotUpdate;
                    if (scriptsToRemove) scriptsToRemove.forEach(function(script) {
                        if (script) {
                            var _document$head2;
                            (_document$head2 = document.head) === null || _document$head2 === void 0 || _document$head2.removeChild(script);
                        }
                    });
                    return [
                        7
                    ];
                case 5:
                    return [
                        2
                    ];
            }
        });
    });
    return _hmrApplyUpdates.apply(this, arguments);
}
function hmrApply(bundle /*: ParcelRequire */ , asset /*:  HMRAsset */ ) {
    var modules = bundle.modules;
    if (!modules) return;
    if (asset.type === "css") reloadCSS();
    else if (asset.type === "js") {
        var deps = asset.depsByBundle[bundle.HMR_BUNDLE_ID];
        if (deps) {
            if (modules[asset.id]) {
                // Remove dependencies that are removed and will become orphaned.
                // This is necessary so that if the asset is added back again, the cache is gone, and we prevent a full page reload.
                var oldDeps = modules[asset.id][1];
                for(var dep in oldDeps)if (!deps[dep] || deps[dep] !== oldDeps[dep]) {
                    var id = oldDeps[dep];
                    var parents = getParents(module.bundle.root, id);
                    if (parents.length === 1) hmrDelete(module.bundle.root, id);
                }
            }
            if (supportsSourceURL) // Global eval. We would use `new Function` here but browser
            // support for source maps is better with eval.
            (0, eval)(asset.output);
            // $FlowFixMe
            var fn = global.parcelHotUpdate[asset.id];
            modules[asset.id] = [
                fn,
                deps
            ];
        } else if (bundle.parent) hmrApply(bundle.parent, asset);
    }
}
function hmrDelete(bundle, id) {
    var modules = bundle.modules;
    if (!modules) return;
    if (modules[id]) {
        // Collect dependencies that will become orphaned when this module is deleted.
        var deps = modules[id][1];
        var orphans = [];
        for(var dep in deps){
            var parents = getParents(module.bundle.root, deps[dep]);
            if (parents.length === 1) orphans.push(deps[dep]);
        }
        // Delete the module. This must be done before deleting dependencies in case of circular dependencies.
        delete modules[id];
        delete bundle.cache[id];
        // Now delete the orphans.
        orphans.forEach(function(id) {
            hmrDelete(module.bundle.root, id);
        });
    } else if (bundle.parent) hmrDelete(bundle.parent, id);
}
function hmrAcceptCheck(bundle /*: ParcelRequire */ , id /*: string */ , depsByBundle /*: ?{ [string]: { [string]: string } }*/ ) {
    if (hmrAcceptCheckOne(bundle, id, depsByBundle)) return true;
    // Traverse parents breadth first. All possible ancestries must accept the HMR update, or we'll reload.
    var parents = getParents(module.bundle.root, id);
    var accepted = false;
    while(parents.length > 0){
        var v = parents.shift();
        var a = hmrAcceptCheckOne(v[0], v[1], null);
        if (a) // If this parent accepts, stop traversing upward, but still consider siblings.
        accepted = true;
        else {
            var _parents;
            // Otherwise, queue the parents in the next level upward.
            var p = getParents(module.bundle.root, v[1]);
            if (p.length === 0) {
                // If there are no parents, then we've reached an entry without accepting. Reload.
                accepted = false;
                break;
            }
            (_parents = parents).push.apply(_parents, (0, _toConsumableArray._)(p));
        }
    }
    return accepted;
}
function hmrAcceptCheckOne(bundle /*: ParcelRequire */ , id /*: string */ , depsByBundle /*: ?{ [string]: { [string]: string } }*/ ) {
    var modules = bundle.modules;
    if (!modules) return;
    if (depsByBundle && !depsByBundle[bundle.HMR_BUNDLE_ID]) {
        // If we reached the root bundle without finding where the asset should go,
        // there's nothing to do. Mark as "accepted" so we don't reload the page.
        if (!bundle.parent) return true;
        return hmrAcceptCheck(bundle.parent, id, depsByBundle);
    }
    if (checkedAssets[id]) return true;
    checkedAssets[id] = true;
    var cached = bundle.cache[id];
    assetsToDispose.push([
        bundle,
        id
    ]);
    if (!cached || cached.hot && cached.hot._acceptCallbacks.length) {
        assetsToAccept.push([
            bundle,
            id
        ]);
        return true;
    }
}
function hmrDispose(bundle /*: ParcelRequire */ , id /*: string */ ) {
    var cached = bundle.cache[id];
    bundle.hotData[id] = {};
    if (cached && cached.hot) cached.hot.data = bundle.hotData[id];
    if (cached && cached.hot && cached.hot._disposeCallbacks.length) cached.hot._disposeCallbacks.forEach(function(cb) {
        cb(bundle.hotData[id]);
    });
    delete bundle.cache[id];
}
function hmrAccept(bundle /*: ParcelRequire */ , id /*: string */ ) {
    // Execute the module.
    bundle(id);
    // Run the accept callbacks in the new version of the module.
    var cached = bundle.cache[id];
    if (cached && cached.hot && cached.hot._acceptCallbacks.length) cached.hot._acceptCallbacks.forEach(function(cb) {
        var assetsToAlsoAccept = cb(function() {
            return getParents(module.bundle.root, id);
        });
        if (assetsToAlsoAccept && assetsToAccept.length) {
            assetsToAlsoAccept.forEach(function(a) {
                hmrDispose(a[0], a[1]);
            });
            // $FlowFixMe[method-unbinding]
            assetsToAccept.push.apply(assetsToAccept, assetsToAlsoAccept);
        }
    });
}

},{"@swc/helpers/_/_async_to_generator":"jWKNZ","@swc/helpers/_/_to_consumable_array":"79eRw","@swc/helpers/_/_ts_generator":"al8kL"}],"jWKNZ":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_async_to_generator", function() {
    return _async_to_generator;
});
parcelHelpers.export(exports, "_", function() {
    return _async_to_generator;
});
function asyncGeneratorStep(gen, resolve, reject, _next, _throw, key, arg) {
    try {
        var info = gen[key](arg);
        var value = info.value;
    } catch (error) {
        reject(error);
        return;
    }
    if (info.done) resolve(value);
    else Promise.resolve(value).then(_next, _throw);
}
function _async_to_generator(fn) {
    return function() {
        var self = this, args = arguments;
        return new Promise(function(resolve, reject) {
            var gen = fn.apply(self, args);
            function _next(value) {
                asyncGeneratorStep(gen, resolve, reject, _next, _throw, "next", value);
            }
            function _throw(err) {
                asyncGeneratorStep(gen, resolve, reject, _next, _throw, "throw", err);
            }
            _next(undefined);
        });
    };
}

},{"@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"3Qwoy":[function(require,module,exports) {
exports.interopDefault = function(a) {
    return a && a.__esModule ? a : {
        default: a
    };
};
exports.defineInteropFlag = function(a) {
    Object.defineProperty(a, "__esModule", {
        value: true
    });
};
exports.exportAll = function(source, dest) {
    Object.keys(source).forEach(function(key) {
        if (key === "default" || key === "__esModule" || Object.prototype.hasOwnProperty.call(dest, key)) return;
        Object.defineProperty(dest, key, {
            enumerable: true,
            get: function get() {
                return source[key];
            }
        });
    });
    return dest;
};
exports.export = function(dest, destName, get) {
    Object.defineProperty(dest, destName, {
        enumerable: true,
        get: get
    });
};

},{}],"79eRw":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_to_consumable_array", function() {
    return _to_consumable_array;
});
parcelHelpers.export(exports, "_", function() {
    return _to_consumable_array;
});
var _arrayWithoutHolesJs = require("./_array_without_holes.js");
var _iterableToArrayJs = require("./_iterable_to_array.js");
var _nonIterableSpreadJs = require("./_non_iterable_spread.js");
var _unsupportedIterableToArrayJs = require("./_unsupported_iterable_to_array.js");
function _to_consumable_array(arr) {
    return (0, _arrayWithoutHolesJs._array_without_holes)(arr) || (0, _iterableToArrayJs._iterable_to_array)(arr) || (0, _unsupportedIterableToArrayJs._unsupported_iterable_to_array)(arr) || (0, _nonIterableSpreadJs._non_iterable_spread)();
}

},{"./_array_without_holes.js":"3QRUb","./_iterable_to_array.js":"5HFUj","./_non_iterable_spread.js":"25UZY","./_unsupported_iterable_to_array.js":"aofca","@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"3QRUb":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_array_without_holes", function() {
    return _array_without_holes;
});
parcelHelpers.export(exports, "_", function() {
    return _array_without_holes;
});
var _arrayLikeToArrayJs = require("./_array_like_to_array.js");
function _array_without_holes(arr) {
    if (Array.isArray(arr)) return (0, _arrayLikeToArrayJs._array_like_to_array)(arr);
}

},{"./_array_like_to_array.js":"8pQhE","@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"8pQhE":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_array_like_to_array", function() {
    return _array_like_to_array;
});
parcelHelpers.export(exports, "_", function() {
    return _array_like_to_array;
});
function _array_like_to_array(arr, len) {
    if (len == null || len > arr.length) len = arr.length;
    for(var i = 0, arr2 = new Array(len); i < len; i++)arr2[i] = arr[i];
    return arr2;
}

},{"@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"5HFUj":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_iterable_to_array", function() {
    return _iterable_to_array;
});
parcelHelpers.export(exports, "_", function() {
    return _iterable_to_array;
});
function _iterable_to_array(iter) {
    if (typeof Symbol !== "undefined" && iter[Symbol.iterator] != null || iter["@@iterator"] != null) return Array.from(iter);
}

},{"@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"25UZY":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_non_iterable_spread", function() {
    return _non_iterable_spread;
});
parcelHelpers.export(exports, "_", function() {
    return _non_iterable_spread;
});
function _non_iterable_spread() {
    throw new TypeError("Invalid attempt to spread non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method.");
}

},{"@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"aofca":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_unsupported_iterable_to_array", function() {
    return _unsupported_iterable_to_array;
});
parcelHelpers.export(exports, "_", function() {
    return _unsupported_iterable_to_array;
});
var _arrayLikeToArrayJs = require("./_array_like_to_array.js");
function _unsupported_iterable_to_array(o, minLen) {
    if (!o) return;
    if (typeof o === "string") return (0, _arrayLikeToArrayJs._array_like_to_array)(o, minLen);
    var n = Object.prototype.toString.call(o).slice(8, -1);
    if (n === "Object" && o.constructor) n = o.constructor.name;
    if (n === "Map" || n === "Set") return Array.from(n);
    if (n === "Arguments" || /^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)) return (0, _arrayLikeToArrayJs._array_like_to_array)(o, minLen);
}

},{"./_array_like_to_array.js":"8pQhE","@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"al8kL":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_", function() {
    return 0, _tslib.__generator;
});
parcelHelpers.export(exports, "_ts_generator", function() {
    return 0, _tslib.__generator;
});
var _tslib = require("tslib");

},{"tslib":"lMezb","@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"lMezb":[function(require,module,exports) {
/******************************************************************************
Copyright (c) Microsoft Corporation.

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY
AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
PERFORMANCE OF THIS SOFTWARE.
***************************************************************************** */ /* global Reflect, Promise, SuppressedError, Symbol */ var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "__extends", function() {
    return __extends;
});
parcelHelpers.export(exports, "__assign", function() {
    return __assign;
});
parcelHelpers.export(exports, "__rest", function() {
    return __rest;
});
parcelHelpers.export(exports, "__decorate", function() {
    return __decorate;
});
parcelHelpers.export(exports, "__param", function() {
    return __param;
});
parcelHelpers.export(exports, "__esDecorate", function() {
    return __esDecorate;
});
parcelHelpers.export(exports, "__runInitializers", function() {
    return __runInitializers;
});
parcelHelpers.export(exports, "__propKey", function() {
    return __propKey;
});
parcelHelpers.export(exports, "__setFunctionName", function() {
    return __setFunctionName;
});
parcelHelpers.export(exports, "__metadata", function() {
    return __metadata;
});
parcelHelpers.export(exports, "__awaiter", function() {
    return __awaiter;
});
parcelHelpers.export(exports, "__generator", function() {
    return __generator;
});
parcelHelpers.export(exports, "__createBinding", function() {
    return __createBinding;
});
parcelHelpers.export(exports, "__exportStar", function() {
    return __exportStar;
});
parcelHelpers.export(exports, "__values", function() {
    return __values;
});
parcelHelpers.export(exports, "__read", function() {
    return __read;
});
/** @deprecated */ parcelHelpers.export(exports, "__spread", function() {
    return __spread;
});
/** @deprecated */ parcelHelpers.export(exports, "__spreadArrays", function() {
    return __spreadArrays;
});
parcelHelpers.export(exports, "__spreadArray", function() {
    return __spreadArray;
});
parcelHelpers.export(exports, "__await", function() {
    return __await;
});
parcelHelpers.export(exports, "__asyncGenerator", function() {
    return __asyncGenerator;
});
parcelHelpers.export(exports, "__asyncDelegator", function() {
    return __asyncDelegator;
});
parcelHelpers.export(exports, "__asyncValues", function() {
    return __asyncValues;
});
parcelHelpers.export(exports, "__makeTemplateObject", function() {
    return __makeTemplateObject;
});
parcelHelpers.export(exports, "__importStar", function() {
    return __importStar;
});
parcelHelpers.export(exports, "__importDefault", function() {
    return __importDefault;
});
parcelHelpers.export(exports, "__classPrivateFieldGet", function() {
    return __classPrivateFieldGet;
});
parcelHelpers.export(exports, "__classPrivateFieldSet", function() {
    return __classPrivateFieldSet;
});
parcelHelpers.export(exports, "__classPrivateFieldIn", function() {
    return __classPrivateFieldIn;
});
parcelHelpers.export(exports, "__addDisposableResource", function() {
    return __addDisposableResource;
});
parcelHelpers.export(exports, "__disposeResources", function() {
    return __disposeResources;
});
var _typeOf = require("@swc/helpers/_/_type_of");
var extendStatics = function extendStatics1(d, b) {
    extendStatics = Object.setPrototypeOf || ({
        __proto__: []
    }) instanceof Array && function(d, b) {
        d.__proto__ = b;
    } || function(d, b) {
        for(var p in b)if (Object.prototype.hasOwnProperty.call(b, p)) d[p] = b[p];
    };
    return extendStatics(d, b);
};
function __extends(d, b) {
    if (typeof b !== "function" && b !== null) throw new TypeError("Class extends value " + String(b) + " is not a constructor or null");
    extendStatics(d, b);
    function __() {
        this.constructor = d;
    }
    d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
}
var __assign = function __assign1() {
    __assign = Object.assign || function __assign(t) {
        for(var s, i = 1, n = arguments.length; i < n; i++){
            s = arguments[i];
            for(var p in s)if (Object.prototype.hasOwnProperty.call(s, p)) t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
function __rest(s, e) {
    var t = {};
    for(var p in s)if (Object.prototype.hasOwnProperty.call(s, p) && e.indexOf(p) < 0) t[p] = s[p];
    if (s != null && typeof Object.getOwnPropertySymbols === "function") {
        for(var i = 0, p = Object.getOwnPropertySymbols(s); i < p.length; i++)if (e.indexOf(p[i]) < 0 && Object.prototype.propertyIsEnumerable.call(s, p[i])) t[p[i]] = s[p[i]];
    }
    return t;
}
function __decorate(decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for(var i = decorators.length - 1; i >= 0; i--)if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
}
function __param(paramIndex, decorator) {
    return function(target, key) {
        decorator(target, key, paramIndex);
    };
}
function __esDecorate(ctor, descriptorIn, decorators, contextIn, initializers, extraInitializers) {
    function accept(f) {
        if (f !== void 0 && typeof f !== "function") throw new TypeError("Function expected");
        return f;
    }
    var kind = contextIn.kind, key = kind === "getter" ? "get" : kind === "setter" ? "set" : "value";
    var target = !descriptorIn && ctor ? contextIn["static"] ? ctor : ctor.prototype : null;
    var descriptor = descriptorIn || (target ? Object.getOwnPropertyDescriptor(target, contextIn.name) : {});
    var _, done = false;
    for(var i = decorators.length - 1; i >= 0; i--){
        var context = {};
        for(var p in contextIn)context[p] = p === "access" ? {} : contextIn[p];
        for(var p in contextIn.access)context.access[p] = contextIn.access[p];
        context.addInitializer = function(f) {
            if (done) throw new TypeError("Cannot add initializers after decoration has completed");
            extraInitializers.push(accept(f || null));
        };
        var result = (0, decorators[i])(kind === "accessor" ? {
            get: descriptor.get,
            set: descriptor.set
        } : descriptor[key], context);
        if (kind === "accessor") {
            if (result === void 0) continue;
            if (result === null || typeof result !== "object") throw new TypeError("Object expected");
            if (_ = accept(result.get)) descriptor.get = _;
            if (_ = accept(result.set)) descriptor.set = _;
            if (_ = accept(result.init)) initializers.unshift(_);
        } else if (_ = accept(result)) {
            if (kind === "field") initializers.unshift(_);
            else descriptor[key] = _;
        }
    }
    if (target) Object.defineProperty(target, contextIn.name, descriptor);
    done = true;
}
function __runInitializers(thisArg, initializers, value) {
    var useValue = arguments.length > 2;
    for(var i = 0; i < initializers.length; i++)value = useValue ? initializers[i].call(thisArg, value) : initializers[i].call(thisArg);
    return useValue ? value : void 0;
}
function __propKey(x) {
    return (typeof x === "undefined" ? "undefined" : (0, _typeOf._)(x)) === "symbol" ? x : "".concat(x);
}
function __setFunctionName(f, name, prefix) {
    if ((typeof name === "undefined" ? "undefined" : (0, _typeOf._)(name)) === "symbol") name = name.description ? "[".concat(name.description, "]") : "";
    return Object.defineProperty(f, "name", {
        configurable: true,
        value: prefix ? "".concat(prefix, " ", name) : name
    });
}
function __metadata(metadataKey, metadataValue) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(metadataKey, metadataValue);
}
function __awaiter(thisArg, _arguments, P, generator) {
    function adopt(value) {
        return value instanceof P ? value : new P(function(resolve) {
            resolve(value);
        });
    }
    return new (P || (P = Promise))(function(resolve, reject) {
        function fulfilled(value) {
            try {
                step(generator.next(value));
            } catch (e) {
                reject(e);
            }
        }
        function rejected(value) {
            try {
                step(generator["throw"](value));
            } catch (e) {
                reject(e);
            }
        }
        function step(result) {
            result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected);
        }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
}
function __generator(thisArg, body) {
    var _ = {
        label: 0,
        sent: function sent() {
            if (t[0] & 1) throw t[1];
            return t[1];
        },
        trys: [],
        ops: []
    }, f, y, t, g;
    return g = {
        next: verb(0),
        "throw": verb(1),
        "return": verb(2)
    }, typeof Symbol === "function" && (g[Symbol.iterator] = function() {
        return this;
    }), g;
    function verb(n) {
        return function(v) {
            return step([
                n,
                v
            ]);
        };
    }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while(g && (g = 0, op[0] && (_ = 0)), _)try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [
                op[0] & 2,
                t.value
            ];
            switch(op[0]){
                case 0:
                case 1:
                    t = op;
                    break;
                case 4:
                    _.label++;
                    return {
                        value: op[1],
                        done: false
                    };
                case 5:
                    _.label++;
                    y = op[1];
                    op = [
                        0
                    ];
                    continue;
                case 7:
                    op = _.ops.pop();
                    _.trys.pop();
                    continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) {
                        _ = 0;
                        continue;
                    }
                    if (op[0] === 3 && (!t || op[1] > t[0] && op[1] < t[3])) {
                        _.label = op[1];
                        break;
                    }
                    if (op[0] === 6 && _.label < t[1]) {
                        _.label = t[1];
                        t = op;
                        break;
                    }
                    if (t && _.label < t[2]) {
                        _.label = t[2];
                        _.ops.push(op);
                        break;
                    }
                    if (t[2]) _.ops.pop();
                    _.trys.pop();
                    continue;
            }
            op = body.call(thisArg, _);
        } catch (e) {
            op = [
                6,
                e
            ];
            y = 0;
        } finally{
            f = t = 0;
        }
        if (op[0] & 5) throw op[1];
        return {
            value: op[0] ? op[1] : void 0,
            done: true
        };
    }
}
var __createBinding = Object.create ? function __createBinding(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) desc = {
        enumerable: true,
        get: function get() {
            return m[k];
        }
    };
    Object.defineProperty(o, k2, desc);
} : function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
};
function __exportStar(m, o) {
    for(var p in m)if (p !== "default" && !Object.prototype.hasOwnProperty.call(o, p)) __createBinding(o, m, p);
}
function __values(o) {
    var s = typeof Symbol === "function" && Symbol.iterator, m = s && o[s], i = 0;
    if (m) return m.call(o);
    if (o && typeof o.length === "number") return {
        next: function next() {
            if (o && i >= o.length) o = void 0;
            return {
                value: o && o[i++],
                done: !o
            };
        }
    };
    throw new TypeError(s ? "Object is not iterable." : "Symbol.iterator is not defined.");
}
function __read(o, n) {
    var m = typeof Symbol === "function" && o[Symbol.iterator];
    if (!m) return o;
    var i = m.call(o), r, ar = [], e;
    try {
        while((n === void 0 || n-- > 0) && !(r = i.next()).done)ar.push(r.value);
    } catch (error) {
        e = {
            error: error
        };
    } finally{
        try {
            if (r && !r.done && (m = i["return"])) m.call(i);
        } finally{
            if (e) throw e.error;
        }
    }
    return ar;
}
function __spread() {
    for(var ar = [], i = 0; i < arguments.length; i++)ar = ar.concat(__read(arguments[i]));
    return ar;
}
function __spreadArrays() {
    for(var s = 0, i = 0, il = arguments.length; i < il; i++)s += arguments[i].length;
    for(var r = Array(s), k = 0, i = 0; i < il; i++)for(var a = arguments[i], j = 0, jl = a.length; j < jl; j++, k++)r[k] = a[j];
    return r;
}
function __spreadArray(to, from, pack) {
    if (pack || arguments.length === 2) {
        for(var i = 0, l = from.length, ar; i < l; i++)if (ar || !(i in from)) {
            if (!ar) ar = Array.prototype.slice.call(from, 0, i);
            ar[i] = from[i];
        }
    }
    return to.concat(ar || Array.prototype.slice.call(from));
}
function __await(v) {
    return this instanceof __await ? (this.v = v, this) : new __await(v);
}
function __asyncGenerator(thisArg, _arguments, generator) {
    if (!Symbol.asyncIterator) throw new TypeError("Symbol.asyncIterator is not defined.");
    var g = generator.apply(thisArg, _arguments || []), i, q = [];
    return i = {}, verb("next"), verb("throw"), verb("return", awaitReturn), i[Symbol.asyncIterator] = function() {
        return this;
    }, i;
    function awaitReturn(f) {
        return function(v) {
            return Promise.resolve(v).then(f, reject);
        };
    }
    function verb(n, f) {
        if (g[n]) {
            i[n] = function(v) {
                return new Promise(function(a, b) {
                    q.push([
                        n,
                        v,
                        a,
                        b
                    ]) > 1 || resume(n, v);
                });
            };
            if (f) i[n] = f(i[n]);
        }
    }
    function resume(n, v) {
        try {
            step(g[n](v));
        } catch (e) {
            settle(q[0][3], e);
        }
    }
    function step(r) {
        r.value instanceof __await ? Promise.resolve(r.value.v).then(fulfill, reject) : settle(q[0][2], r);
    }
    function fulfill(value) {
        resume("next", value);
    }
    function reject(value) {
        resume("throw", value);
    }
    function settle(f, v) {
        if (f(v), q.shift(), q.length) resume(q[0][0], q[0][1]);
    }
}
function __asyncDelegator(o) {
    var i, p;
    return i = {}, verb("next"), verb("throw", function(e) {
        throw e;
    }), verb("return"), i[Symbol.iterator] = function() {
        return this;
    }, i;
    function verb(n, f) {
        i[n] = o[n] ? function(v) {
            return (p = !p) ? {
                value: __await(o[n](v)),
                done: false
            } : f ? f(v) : v;
        } : f;
    }
}
function __asyncValues(o) {
    if (!Symbol.asyncIterator) throw new TypeError("Symbol.asyncIterator is not defined.");
    var m = o[Symbol.asyncIterator], i;
    return m ? m.call(o) : (o = typeof __values === "function" ? __values(o) : o[Symbol.iterator](), i = {}, verb("next"), verb("throw"), verb("return"), i[Symbol.asyncIterator] = function() {
        return this;
    }, i);
    function verb(n) {
        i[n] = o[n] && function(v) {
            return new Promise(function(resolve, reject) {
                v = o[n](v), settle(resolve, reject, v.done, v.value);
            });
        };
    }
    function settle(resolve, reject, d, v) {
        Promise.resolve(v).then(function(v) {
            resolve({
                value: v,
                done: d
            });
        }, reject);
    }
}
function __makeTemplateObject(cooked, raw) {
    if (Object.defineProperty) Object.defineProperty(cooked, "raw", {
        value: raw
    });
    else cooked.raw = raw;
    return cooked;
}
var __setModuleDefault = Object.create ? function __setModuleDefault(o, v) {
    Object.defineProperty(o, "default", {
        enumerable: true,
        value: v
    });
} : function(o, v) {
    o["default"] = v;
};
function __importStar(mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) {
        for(var k in mod)if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    }
    __setModuleDefault(result, mod);
    return result;
}
function __importDefault(mod) {
    return mod && mod.__esModule ? mod : {
        default: mod
    };
}
function __classPrivateFieldGet(receiver, state, kind, f) {
    if (kind === "a" && !f) throw new TypeError("Private accessor was defined without a getter");
    if (typeof state === "function" ? receiver !== state || !f : !state.has(receiver)) throw new TypeError("Cannot read private member from an object whose class did not declare it");
    return kind === "m" ? f : kind === "a" ? f.call(receiver) : f ? f.value : state.get(receiver);
}
function __classPrivateFieldSet(receiver, state, value, kind, f) {
    if (kind === "m") throw new TypeError("Private method is not writable");
    if (kind === "a" && !f) throw new TypeError("Private accessor was defined without a setter");
    if (typeof state === "function" ? receiver !== state || !f : !state.has(receiver)) throw new TypeError("Cannot write private member to an object whose class did not declare it");
    return kind === "a" ? f.call(receiver, value) : f ? f.value = value : state.set(receiver, value), value;
}
function __classPrivateFieldIn(state, receiver) {
    if (receiver === null || typeof receiver !== "object" && typeof receiver !== "function") throw new TypeError("Cannot use 'in' operator on non-object");
    return typeof state === "function" ? receiver === state : state.has(receiver);
}
function __addDisposableResource(env, value, async) {
    if (value !== null && value !== void 0) {
        if (typeof value !== "object" && typeof value !== "function") throw new TypeError("Object expected.");
        var dispose, inner;
        if (async) {
            if (!Symbol.asyncDispose) throw new TypeError("Symbol.asyncDispose is not defined.");
            dispose = value[Symbol.asyncDispose];
        }
        if (dispose === void 0) {
            if (!Symbol.dispose) throw new TypeError("Symbol.dispose is not defined.");
            dispose = value[Symbol.dispose];
            if (async) inner = dispose;
        }
        if (typeof dispose !== "function") throw new TypeError("Object not disposable.");
        if (inner) dispose = function dispose() {
            try {
                inner.call(this);
            } catch (e) {
                return Promise.reject(e);
            }
        };
        env.stack.push({
            value: value,
            dispose: dispose,
            async: async
        });
    } else if (async) env.stack.push({
        async: true
    });
    return value;
}
var _SuppressedError = typeof SuppressedError === "function" ? SuppressedError : function _SuppressedError(error, suppressed, message) {
    var e = new Error(message);
    return e.name = "SuppressedError", e.error = error, e.suppressed = suppressed, e;
};
function __disposeResources(env) {
    function fail(e) {
        env.error = env.hasError ? new _SuppressedError(e, env.error, "An error was suppressed during disposal.") : e;
        env.hasError = true;
    }
    function next() {
        while(env.stack.length){
            var rec = env.stack.pop();
            try {
                var result = rec.dispose && rec.dispose.call(rec.value);
                if (rec.async) return Promise.resolve(result).then(next, function(e) {
                    fail(e);
                    return next();
                });
            } catch (e) {
                fail(e);
            }
        }
        if (env.hasError) throw env.error;
    }
    return next();
}
exports.default = {
    __extends: __extends,
    __assign: __assign,
    __rest: __rest,
    __decorate: __decorate,
    __param: __param,
    __metadata: __metadata,
    __awaiter: __awaiter,
    __generator: __generator,
    __createBinding: __createBinding,
    __exportStar: __exportStar,
    __values: __values,
    __read: __read,
    __spread: __spread,
    __spreadArrays: __spreadArrays,
    __spreadArray: __spreadArray,
    __await: __await,
    __asyncGenerator: __asyncGenerator,
    __asyncDelegator: __asyncDelegator,
    __asyncValues: __asyncValues,
    __makeTemplateObject: __makeTemplateObject,
    __importStar: __importStar,
    __importDefault: __importDefault,
    __classPrivateFieldGet: __classPrivateFieldGet,
    __classPrivateFieldSet: __classPrivateFieldSet,
    __classPrivateFieldIn: __classPrivateFieldIn,
    __addDisposableResource: __addDisposableResource,
    __disposeResources: __disposeResources
};

},{"@swc/helpers/_/_type_of":"lQLoi","@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"lQLoi":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_type_of", function() {
    return _type_of;
});
parcelHelpers.export(exports, "_", function() {
    return _type_of;
});
function _type_of(obj) {
    "@swc/helpers - typeof";
    return obj && typeof Symbol !== "undefined" && obj.constructor === Symbol ? "symbol" : typeof obj;
}

},{"@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"42jTR":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
var _dropzone = require("../src/dropzone");
var _dropzoneDefault = parcelHelpers.interopDefault(_dropzone);
window.Dropzone = (0, _dropzoneDefault.default);
exports.default = (0, _dropzoneDefault.default);

},{"../src/dropzone":"1u0OZ","@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"1u0OZ":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "default", function() {
    return Dropzone;
});
parcelHelpers.export(exports, "Dropzone", function() {
    return Dropzone;
});
var _assertThisInitialized = require("@swc/helpers/_/_assert_this_initialized");
var _classCallCheck = require("@swc/helpers/_/_class_call_check");
var _createClass = require("@swc/helpers/_/_create_class");
var _inherits = require("@swc/helpers/_/_inherits");
var _objectSpread = require("@swc/helpers/_/_object_spread");
var _objectSpreadProps = require("@swc/helpers/_/_object_spread_props");
var _possibleConstructorReturn = require("@swc/helpers/_/_possible_constructor_return");
var _createSuper = require("@swc/helpers/_/_create_super");
var _emitter = require("./emitter");
var _emitterDefault = parcelHelpers.interopDefault(_emitter);
var _options = require("./options");
var _optionsDefault = parcelHelpers.interopDefault(_options);
var Dropzone = /*#__PURE__*/ function(Emitter) {
    "use strict";
    (0, _inherits._)(Dropzone, Emitter);
    var _super = (0, _createSuper._)(Dropzone);
    function Dropzone(el, options) {
        (0, _classCallCheck._)(this, Dropzone);
        var _this;
        _this = _super.call(this);
        var fallback, left;
        _this.element = el;
        _this.clickableElements = [];
        _this.listeners = [];
        _this.files = []; // All files
        if (typeof _this.element === "string") _this.element = document.querySelector(_this.element);
        // make sure we actually have an HTML Element
        if (_this.element === null || !_this.element instanceof HTMLElement) throw new Error("Invalid dropzone element: not an instance of HTMLElement.");
        if (_this.element.dropzone) throw new Error("Dropzone already attached.");
        // Now add this dropzone to the instances.
        Dropzone.instances.push((0, _assertThisInitialized._)(_this));
        // Put the dropzone inside the element itself.
        _this.element.dropzone = (0, _assertThisInitialized._)(_this);
        var elementOptions = (left = Dropzone.optionsForElement(_this.element)) != null ? left : {};
        _this.options = Object.assign({}, (0, _optionsDefault.default), elementOptions, options != null ? options : {});
        _this.options.previewTemplate = _this.options.previewTemplate.replace(/\n*/g, "");
        // If the browser failed, just call the fallback and leave
        if (_this.options.forceFallback || !Dropzone.isBrowserSupported()) return (0, _possibleConstructorReturn._)(_this, _this.options.fallback.call((0, _assertThisInitialized._)(_this)));
        // @options.url = @element.getAttribute "action" unless @options.url?
        if (_this.options.url == null) _this.options.url = _this.element.getAttribute("action");
        if (!_this.options.url) throw new Error("No URL provided.");
        if (_this.options.uploadMultiple && _this.options.chunking) throw new Error("You cannot set both: uploadMultiple and chunking.");
        if (_this.options.binaryBody && _this.options.uploadMultiple) throw new Error("You cannot set both: binaryBody and uploadMultiple.");
        if (typeof _this.options.method === "string") _this.options.method = _this.options.method.toUpperCase();
        if ((fallback = _this.getExistingFallback()) && fallback.parentNode) // Remove the fallback
        fallback.parentNode.removeChild(fallback);
        // Display previews in the previewsContainer element or the Dropzone element unless explicitly set to false
        if (_this.options.previewsContainer !== false) {
            if (_this.options.previewsContainer) _this.previewsContainer = Dropzone.getElement(_this.options.previewsContainer, "previewsContainer");
            else _this.previewsContainer = _this.element;
        }
        if (_this.options.clickable) {
            if (_this.options.clickable === true) _this.clickableElements = [
                _this.element
            ];
            else _this.clickableElements = Dropzone.getElements(_this.options.clickable, "clickable");
        }
        _this.init();
        return _this;
    }
    (0, _createClass._)(Dropzone, [
        {
            // Returns all files that have been accepted
            key: "getAcceptedFiles",
            value: function getAcceptedFiles() {
                return this.files.filter(function(file) {
                    return file.accepted;
                }).map(function(file) {
                    return file;
                });
            }
        },
        {
            // Returns all files that have been rejected
            // Not sure when that's going to be useful, but added for completeness.
            key: "getRejectedFiles",
            value: function getRejectedFiles() {
                return this.files.filter(function(file) {
                    return !file.accepted;
                }).map(function(file) {
                    return file;
                });
            }
        },
        {
            key: "getFilesWithStatus",
            value: function getFilesWithStatus(status) {
                return this.files.filter(function(file) {
                    return file.status === status;
                }).map(function(file) {
                    return file;
                });
            }
        },
        {
            // Returns all files that are in the queue
            key: "getQueuedFiles",
            value: function getQueuedFiles() {
                return this.getFilesWithStatus(Dropzone.QUEUED);
            }
        },
        {
            key: "getUploadingFiles",
            value: function getUploadingFiles() {
                return this.getFilesWithStatus(Dropzone.UPLOADING);
            }
        },
        {
            key: "getAddedFiles",
            value: function getAddedFiles() {
                return this.getFilesWithStatus(Dropzone.ADDED);
            }
        },
        {
            // Files that are either queued or uploading
            key: "getActiveFiles",
            value: function getActiveFiles() {
                return this.files.filter(function(file) {
                    return file.status === Dropzone.UPLOADING || file.status === Dropzone.QUEUED;
                }).map(function(file) {
                    return file;
                });
            }
        },
        {
            // The function that gets called when Dropzone is initialized. You
            // can (and should) setup event listeners inside this function.
            key: "init",
            value: function init() {
                var _this = this;
                // In case it isn't set already
                if (this.element.tagName === "form") this.element.setAttribute("enctype", "multipart/form-data");
                if (this.element.classList.contains("dropzone") && !this.element.querySelector(".dz-message")) this.element.appendChild(Dropzone.createElement('<div class="dz-default dz-message"><button class="dz-button" type="button">'.concat(this.options.dictDefaultMessage, "</button></div>")));
                if (this.clickableElements.length) {
                    var setupHiddenFileInput = function() {
                        var _this_hiddenFileInput_parentNode;
                        if (_this.hiddenFileInput) (_this_hiddenFileInput_parentNode = _this.hiddenFileInput.parentNode) === null || _this_hiddenFileInput_parentNode === void 0 ? void 0 : _this_hiddenFileInput_parentNode.removeChild(_this.hiddenFileInput);
                        _this.hiddenFileInput = document.createElement("input");
                        _this.hiddenFileInput.setAttribute("type", "file");
                        _this.hiddenFileInput.setAttribute("form", _this.element.id);
                        if (_this.options.maxFiles === null || _this.options.maxFiles > 1) _this.hiddenFileInput.setAttribute("multiple", "multiple");
                        _this.hiddenFileInput.className = "dz-hidden-input";
                        if (_this.options.acceptedFiles !== null) _this.hiddenFileInput.setAttribute("accept", _this.options.acceptedFiles);
                        if (_this.options.capture !== null) _this.hiddenFileInput.setAttribute("capture", _this.options.capture);
                        // Making sure that no one can "tab" into this field.
                        _this.hiddenFileInput.setAttribute("tabindex", "-1");
                        // Add arialabel for a11y
                        _this.hiddenFileInput.setAttribute("aria-label", "dropzone hidden input");
                        // Not setting `display="none"` because some browsers don't accept clicks
                        // on elements that aren't displayed.
                        _this.hiddenFileInput.style.visibility = "hidden";
                        _this.hiddenFileInput.style.position = "absolute";
                        _this.hiddenFileInput.style.top = "0";
                        _this.hiddenFileInput.style.left = "0";
                        _this.hiddenFileInput.style.height = "0";
                        _this.hiddenFileInput.style.width = "0";
                        Dropzone.getElement(_this.options.hiddenInputContainer, "hiddenInputContainer").appendChild(_this.hiddenFileInput);
                        _this.hiddenFileInput.addEventListener("change", function() {
                            var files = _this.hiddenFileInput.files;
                            var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                            if (files.length) try {
                                for(var _iterator = files[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                                    var file = _step.value;
                                    _this.addFile(file);
                                }
                            } catch (err) {
                                _didIteratorError = true;
                                _iteratorError = err;
                            } finally{
                                try {
                                    if (!_iteratorNormalCompletion && _iterator.return != null) {
                                        _iterator.return();
                                    }
                                } finally{
                                    if (_didIteratorError) {
                                        throw _iteratorError;
                                    }
                                }
                            }
                            _this.emit("addedfiles", files);
                            setupHiddenFileInput();
                        });
                    };
                    setupHiddenFileInput();
                }
                this.URL = window.URL !== null ? window.URL : window.webkitURL;
                var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                try {
                    // Setup all event listeners on the Dropzone object itself.
                    // They're not in @setupEventListeners() because they shouldn't be removed
                    // again when the dropzone gets disabled.
                    for(var _iterator = this.events[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                        var eventName = _step.value;
                        this.on(eventName, this.options[eventName]);
                    }
                } catch (err) {
                    _didIteratorError = true;
                    _iteratorError = err;
                } finally{
                    try {
                        if (!_iteratorNormalCompletion && _iterator.return != null) {
                            _iterator.return();
                        }
                    } finally{
                        if (_didIteratorError) {
                            throw _iteratorError;
                        }
                    }
                }
                this.on("uploadprogress", function() {
                    return _this.updateTotalUploadProgress();
                });
                this.on("removedfile", function() {
                    return _this.updateTotalUploadProgress();
                });
                this.on("canceled", function(file) {
                    return _this.emit("complete", file);
                });
                // Emit a `queuecomplete` event if all files finished uploading.
                this.on("complete", function(file) {
                    if (_this.getAddedFiles().length === 0 && _this.getUploadingFiles().length === 0 && _this.getQueuedFiles().length === 0) // This needs to be deferred so that `queuecomplete` really triggers after `complete`
                    return setTimeout(function() {
                        return _this.emit("queuecomplete");
                    }, 0);
                });
                var containsFiles = function containsFiles(e) {
                    return e.dataTransfer.types && e.dataTransfer.types.includes("Files");
                };
                var noPropagation = function noPropagation(e) {
                    // If there are no files, we don't want to stop
                    // propagation so we don't interfere with other
                    // drag and drop behaviour.
                    if (!containsFiles(e)) return;
                    e.stopPropagation();
                    return e.preventDefault();
                };
                // Create the listeners
                this.listeners = [
                    {
                        element: this.element,
                        events: {
                            dragstart: function(e) {
                                return _this.emit("dragstart", e);
                            },
                            dragenter: function(e) {
                                noPropagation(e);
                                return _this.emit("dragenter", e);
                            },
                            dragover: function(e) {
                                // Makes it possible to drag files from chrome's download bar
                                // http://stackoverflow.com/questions/19526430/drag-and-drop-file-uploads-from-chrome-downloads-bar
                                var efct = e.dataTransfer.effectAllowed;
                                e.dataTransfer.dropEffect = "move" === efct || "linkMove" === efct ? "move" : "copy";
                                noPropagation(e);
                                return _this.emit("dragover", e);
                            },
                            dragleave: function(e) {
                                return _this.emit("dragleave", e);
                            },
                            drop: function(e) {
                                noPropagation(e);
                                return _this.drop(e);
                            },
                            dragend: function(e) {
                                return _this.emit("dragend", e);
                            }
                        }
                    }
                ];
                this.clickableElements.forEach(function(clickableElement) {
                    return _this.listeners.push({
                        element: clickableElement,
                        events: {
                            click: function(evt) {
                                // Only the actual dropzone or the message element should trigger file selection
                                if (clickableElement !== _this.element || evt.target === _this.element || Dropzone.elementInside(evt.target, _this.element.querySelector(".dz-message"))) _this.hiddenFileInput.click(); // Forward the click
                                return true;
                            }
                        }
                    });
                });
                this.enable();
                return this.options.init.call(this);
            }
        },
        {
            // Not fully tested yet
            key: "destroy",
            value: function destroy() {
                this.disable();
                this.removeAllFiles(true);
                if (this.hiddenFileInput != null ? this.hiddenFileInput.parentNode : undefined) {
                    this.hiddenFileInput.parentNode.removeChild(this.hiddenFileInput);
                    this.hiddenFileInput = null;
                }
                delete this.element.dropzone;
                return Dropzone.instances.splice(Dropzone.instances.indexOf(this), 1);
            }
        },
        {
            key: "updateTotalUploadProgress",
            value: function updateTotalUploadProgress() {
                var totalUploadProgress;
                var totalBytesSent = 0;
                var totalBytes = 0;
                var activeFiles = this.getActiveFiles();
                if (activeFiles.length) {
                    var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                    try {
                        for(var _iterator = this.getActiveFiles()[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                            var file = _step.value;
                            totalBytesSent += file.upload.bytesSent;
                            totalBytes += file.upload.total;
                        }
                    } catch (err) {
                        _didIteratorError = true;
                        _iteratorError = err;
                    } finally{
                        try {
                            if (!_iteratorNormalCompletion && _iterator.return != null) {
                                _iterator.return();
                            }
                        } finally{
                            if (_didIteratorError) {
                                throw _iteratorError;
                            }
                        }
                    }
                    totalUploadProgress = 100 * totalBytesSent / totalBytes;
                } else totalUploadProgress = 100;
                return this.emit("totaluploadprogress", totalUploadProgress, totalBytes, totalBytesSent);
            }
        },
        {
            // @options.paramName can be a function taking one parameter rather than a string.
            // A parameter name for a file is obtained simply by calling this with an index number.
            key: "_getParamName",
            value: function _getParamName(n) {
                if (typeof this.options.paramName === "function") return this.options.paramName(n);
                else return "".concat(this.options.paramName).concat(this.options.uploadMultiple ? "[".concat(n, "]") : "");
            }
        },
        {
            // If @options.renameFile is a function,
            // the function will be used to rename the file.name before appending it to the formData.
            // MacOS 14+ screenshots contain narrow non-breaking space (U+202F) characters in filenames 
            // (e.g., "Screenshot 2024-01-30 at 10.32.07 AM.png" where the space after "07" and before "AM" is U+202F).
            // This function now replaces these with regular spaces to prevent upload issues and maintain compatibility with MacOS
            key: "_renameFile",
            value: function _renameFile(file) {
                var cleanFile = (0, _objectSpreadProps._)((0, _objectSpread._)({}, file), {
                    name: file.name.replace(/\u202F/g, " ")
                });
                if (typeof this.options.renameFile !== "function") return cleanFile.name;
                return this.options.renameFile(cleanFile);
            }
        },
        {
            // Returns a form that can be used as fallback if the browser does not support DragnDrop
            //
            // If the dropzone is already a form, only the input field and button are returned. Otherwise a complete form element is provided.
            // This code has to pass in IE7 :(
            key: "getFallbackForm",
            value: function getFallbackForm() {
                var existingFallback, form;
                if (existingFallback = this.getExistingFallback()) return existingFallback;
                var fieldsString = '<div class="dz-fallback">';
                if (this.options.dictFallbackText) fieldsString += "<p>".concat(this.options.dictFallbackText, "</p>");
                fieldsString += '<input type="file" name="'.concat(this._getParamName(0), '" ').concat(this.options.uploadMultiple ? 'multiple="multiple"' : undefined, ' /><input type="submit" value="Upload!"></div>');
                var fields = Dropzone.createElement(fieldsString);
                if (this.element.tagName !== "FORM") {
                    form = Dropzone.createElement('<form action="'.concat(this.options.url, '" enctype="multipart/form-data" method="').concat(this.options.method, '"></form>'));
                    form.appendChild(fields);
                } else {
                    // Make sure that the enctype and method attributes are set properly
                    this.element.setAttribute("enctype", "multipart/form-data");
                    this.element.setAttribute("method", this.options.method);
                }
                return form != null ? form : fields;
            }
        },
        {
            // Returns the fallback elements if they exist already
            //
            // This code has to pass in IE7 :(
            key: "getExistingFallback",
            value: function getExistingFallback() {
                var getFallback = function getFallback(elements) {
                    var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                    try {
                        for(var _iterator = elements[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                            var el = _step.value;
                            if (/(^| )fallback($| )/.test(el.className)) return el;
                        }
                    } catch (err) {
                        _didIteratorError = true;
                        _iteratorError = err;
                    } finally{
                        try {
                            if (!_iteratorNormalCompletion && _iterator.return != null) {
                                _iterator.return();
                            }
                        } finally{
                            if (_didIteratorError) {
                                throw _iteratorError;
                            }
                        }
                    }
                };
                for(var _i = 0, _iter = [
                    "div",
                    "form"
                ]; _i < _iter.length; _i++){
                    var tagName = _iter[_i];
                    var fallback;
                    if (fallback = getFallback(this.element.getElementsByTagName(tagName))) return fallback;
                }
            }
        },
        {
            // Activates all listeners stored in @listeners
            key: "setupEventListeners",
            value: function setupEventListeners() {
                return this.listeners.map(function(elementListeners) {
                    return function() {
                        var result = [];
                        for(var event in elementListeners.events){
                            var listener = elementListeners.events[event];
                            result.push(elementListeners.element.addEventListener(event, listener, false));
                        }
                        return result;
                    }();
                });
            }
        },
        {
            // Deactivates all listeners stored in @listeners
            key: "removeEventListeners",
            value: function removeEventListeners() {
                return this.listeners.map(function(elementListeners) {
                    return function() {
                        var result = [];
                        for(var event in elementListeners.events){
                            var listener = elementListeners.events[event];
                            result.push(elementListeners.element.removeEventListener(event, listener, false));
                        }
                        return result;
                    }();
                });
            }
        },
        {
            // Removes all event listeners and cancels all files in the queue or being processed.
            key: "disable",
            value: function disable() {
                var _this = this;
                this.clickableElements.forEach(function(element) {
                    return element.classList.remove("dz-clickable");
                });
                this.removeEventListeners();
                this.disabled = true;
                return this.files.map(function(file) {
                    return _this.cancelUpload(file);
                });
            }
        },
        {
            key: "enable",
            value: function enable() {
                delete this.disabled;
                this.clickableElements.forEach(function(element) {
                    return element.classList.add("dz-clickable");
                });
                return this.setupEventListeners();
            }
        },
        {
            // Returns a nicely formatted filesize
            key: "filesize",
            value: function filesize(size) {
                var selectedSize = 0;
                var selectedUnit = "b";
                if (size > 0) {
                    var units = [
                        "tb",
                        "gb",
                        "mb",
                        "kb",
                        "b"
                    ];
                    for(var i = 0; i < units.length; i++){
                        var unit = units[i];
                        var cutoff = Math.pow(this.options.filesizeBase, 4 - i) / 10;
                        if (size >= cutoff) {
                            selectedSize = size / Math.pow(this.options.filesizeBase, 4 - i);
                            selectedUnit = unit;
                            break;
                        }
                    }
                    selectedSize = Math.round(10 * selectedSize) / 10; // Cutting of digits
                }
                return "<strong>".concat(selectedSize, "</strong> ").concat(this.options.dictFileSizeUnits[selectedUnit]);
            }
        },
        {
            // Adds or removes the `dz-max-files-reached` class from the form.
            key: "_updateMaxFilesReachedClass",
            value: function _updateMaxFilesReachedClass() {
                if (this.options.maxFiles != null && this.getAcceptedFiles().length >= this.options.maxFiles) {
                    if (this.getAcceptedFiles().length === this.options.maxFiles) this.emit("maxfilesreached", this.files);
                    return this.element.classList.add("dz-max-files-reached");
                } else return this.element.classList.remove("dz-max-files-reached");
            }
        },
        {
            key: "drop",
            value: function drop(e) {
                if (!e.dataTransfer) return;
                this.emit("drop", e);
                // Convert the FileList to an Array
                // This is necessary for IE11
                var files = [];
                for(var i = 0; i < e.dataTransfer.files.length; i++)files[i] = e.dataTransfer.files[i];
                // Even if it's a folder, files.length will contain the folders.
                if (files.length) {
                    var items = e.dataTransfer.items;
                    if (items && items.length && items[0].webkitGetAsEntry != null) // The browser supports dropping of folders, so handle items instead of files
                    this._addFilesFromItems(items);
                    else this.handleFiles(files);
                }
                this.emit("addedfiles", files);
            }
        },
        {
            key: "paste",
            value: function paste(e) {
                if (__guard__(e != null ? e.clipboardData : undefined, function(x) {
                    return x.items;
                }) == null) return;
                this.emit("paste", e);
                var items = e.clipboardData.items;
                if (items.length) return this._addFilesFromItems(items);
            }
        },
        {
            key: "handleFiles",
            value: function handleFiles(files) {
                var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                try {
                    for(var _iterator = files[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                        var file = _step.value;
                        this.addFile(file);
                    }
                } catch (err) {
                    _didIteratorError = true;
                    _iteratorError = err;
                } finally{
                    try {
                        if (!_iteratorNormalCompletion && _iterator.return != null) {
                            _iterator.return();
                        }
                    } finally{
                        if (_didIteratorError) {
                            throw _iteratorError;
                        }
                    }
                }
            }
        },
        {
            // When a folder is dropped (or files are pasted), items must be handled
            // instead of files.
            key: "_addFilesFromItems",
            value: function _addFilesFromItems(items) {
                var _this = this;
                return function() {
                    var result = [];
                    var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                    try {
                        for(var _iterator = items[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                            var item = _step.value;
                            var entry;
                            if (item.webkitGetAsEntry != null && (entry = item.webkitGetAsEntry())) {
                                if (entry.isFile) result.push(_this.addFile(item.getAsFile()));
                                else if (entry.isDirectory) // Append all files from that directory to files
                                result.push(_this._addFilesFromDirectory(entry, entry.name));
                                else result.push(undefined);
                            } else if (item.getAsFile != null) {
                                if (item.kind == null || item.kind === "file") result.push(_this.addFile(item.getAsFile()));
                                else result.push(undefined);
                            } else result.push(undefined);
                        }
                    } catch (err) {
                        _didIteratorError = true;
                        _iteratorError = err;
                    } finally{
                        try {
                            if (!_iteratorNormalCompletion && _iterator.return != null) {
                                _iterator.return();
                            }
                        } finally{
                            if (_didIteratorError) {
                                throw _iteratorError;
                            }
                        }
                    }
                    return result;
                }();
            }
        },
        {
            // Goes through the directory, and adds each file it finds recursively
            key: "_addFilesFromDirectory",
            value: function _addFilesFromDirectory(directory, path) {
                var _this = this;
                var dirReader = directory.createReader();
                var errorHandler = function(error) {
                    return __guardMethod__(console, "log", function(o) {
                        return o.log(error);
                    });
                };
                var readEntries = function() {
                    return dirReader.readEntries(function(entries) {
                        if (entries.length > 0) {
                            var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                            try {
                                for(var _iterator = entries[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                                    var entry = _step.value;
                                    if (entry.isFile) entry.file(function(file) {
                                        if (_this.options.ignoreHiddenFiles && file.name.substring(0, 1) === ".") return;
                                        file.fullPath = "".concat(path, "/").concat(file.name);
                                        return _this.addFile(file);
                                    });
                                    else if (entry.isDirectory) _this._addFilesFromDirectory(entry, "".concat(path, "/").concat(entry.name));
                                }
                            } catch (err) {
                                _didIteratorError = true;
                                _iteratorError = err;
                            } finally{
                                try {
                                    if (!_iteratorNormalCompletion && _iterator.return != null) {
                                        _iterator.return();
                                    }
                                } finally{
                                    if (_didIteratorError) {
                                        throw _iteratorError;
                                    }
                                }
                            }
                            // Recursively call readEntries() again, since browser only handle
                            // the first 100 entries.
                            // See: https://developer.mozilla.org/en-US/docs/Web/API/DirectoryReader#readEntries
                            readEntries();
                        }
                        return null;
                    }, errorHandler);
                };
                return readEntries();
            }
        },
        {
            // If `done()` is called without argument the file is accepted
            // If you call it with an error message, the file is rejected
            // (This allows for asynchronous validation)
            //
            // This function checks the filesize, and if the file.type passes the
            // `acceptedFiles` check.
            key: "accept",
            value: function accept(file, done) {
                if (this.options.maxFilesize && file.size > this.options.maxFilesize * 1048576) done(this.options.dictFileTooBig.replace("{{filesize}}", Math.round(file.size / 1024 / 10.24) / 100).replace("{{maxFilesize}}", this.options.maxFilesize));
                else if (!Dropzone.isValidFile(file, this.options.acceptedFiles)) done(this.options.dictInvalidFileType);
                else if (this.options.maxFiles != null && this.getAcceptedFiles().length >= this.options.maxFiles) {
                    done(this.options.dictMaxFilesExceeded.replace("{{maxFiles}}", this.options.maxFiles));
                    this.emit("maxfilesexceeded", file);
                } else this.options.accept.call(this, file, done);
            }
        },
        {
            key: "addFile",
            value: function addFile(file) {
                var _this = this;
                file.upload = {
                    // note: this only works if window.isSecureContext is true, which includes localhost in http
                    uuid: window.isSecureContext ? self.crypto.randomUUID() : Dropzone.uuidv4(),
                    progress: 0,
                    // Setting the total upload size to file.size for the beginning
                    // It's actual different than the size to be transmitted.
                    total: file.size,
                    bytesSent: 0,
                    filename: this._renameFile(file)
                };
                this.files.push(file);
                file.status = Dropzone.ADDED;
                this.emit("addedfile", file);
                this._enqueueThumbnail(file);
                this.accept(file, function(error) {
                    if (error) {
                        file.accepted = false;
                        _this._errorProcessing([
                            file
                        ], error); // Will set the file.status
                    } else {
                        file.accepted = true;
                        if (_this.options.autoQueue) _this.enqueueFile(file);
                         // Will set .accepted = true
                    }
                    _this._updateMaxFilesReachedClass();
                });
            }
        },
        {
            // Wrapper for enqueueFile
            key: "enqueueFiles",
            value: function enqueueFiles(files) {
                var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                try {
                    for(var _iterator = files[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                        var file = _step.value;
                        this.enqueueFile(file);
                    }
                } catch (err) {
                    _didIteratorError = true;
                    _iteratorError = err;
                } finally{
                    try {
                        if (!_iteratorNormalCompletion && _iterator.return != null) {
                            _iterator.return();
                        }
                    } finally{
                        if (_didIteratorError) {
                            throw _iteratorError;
                        }
                    }
                }
                return null;
            }
        },
        {
            key: "enqueueFile",
            value: function enqueueFile(file) {
                var _this = this;
                if (file.status === Dropzone.ADDED && file.accepted === true) {
                    file.status = Dropzone.QUEUED;
                    if (this.options.autoProcessQueue) return setTimeout(function() {
                        return _this.processQueue();
                    }, 0); // Deferring the call
                } else throw new Error("This file can't be queued because it has already been processed or was rejected.");
            }
        },
        {
            key: "_enqueueThumbnail",
            value: function _enqueueThumbnail(file) {
                var _this = this;
                if (this.options.createImageThumbnails && file.type.match(/image.*/) && file.size <= this.options.maxThumbnailFilesize * 1048576) {
                    this._thumbnailQueue.push(file);
                    return setTimeout(function() {
                        return _this._processThumbnailQueue();
                    }, 0); // Deferring the call
                }
            }
        },
        {
            key: "_processThumbnailQueue",
            value: function _processThumbnailQueue() {
                var _this = this;
                if (this._processingThumbnail || this._thumbnailQueue.length === 0) return;
                this._processingThumbnail = true;
                var file = this._thumbnailQueue.shift();
                return this.createThumbnail(file, this.options.thumbnailWidth, this.options.thumbnailHeight, this.options.thumbnailMethod, true, function(dataUrl) {
                    _this.emit("thumbnail", file, dataUrl);
                    _this._processingThumbnail = false;
                    return _this._processThumbnailQueue();
                });
            }
        },
        {
            // Can be called by the user to remove a file
            key: "removeFile",
            value: function removeFile(file) {
                if (file.status === Dropzone.UPLOADING) this.cancelUpload(file);
                this.files = without(this.files, file);
                this.emit("removedfile", file);
                if (this.files.length === 0) return this.emit("reset");
            }
        },
        {
            // Removes all files that aren't currently processed from the list
            key: "removeAllFiles",
            value: function removeAllFiles(cancelIfNecessary) {
                // Create a copy of files since removeFile() changes the @files array.
                if (cancelIfNecessary == null) cancelIfNecessary = false;
                var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                try {
                    for(var _iterator = this.files.slice()[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                        var file = _step.value;
                        if (file.status !== Dropzone.UPLOADING || cancelIfNecessary) this.removeFile(file);
                    }
                } catch (err) {
                    _didIteratorError = true;
                    _iteratorError = err;
                } finally{
                    try {
                        if (!_iteratorNormalCompletion && _iterator.return != null) {
                            _iterator.return();
                        }
                    } finally{
                        if (_didIteratorError) {
                            throw _iteratorError;
                        }
                    }
                }
                return null;
            }
        },
        {
            // Resizes an image before it gets sent to the server. This function is the default behavior of
            // `options.transformFile` if `resizeWidth` or `resizeHeight` are set. The callback is invoked with
            // the resized blob.
            key: "resizeImage",
            value: function resizeImage(file, width, height, resizeMethod, callback) {
                var _this = this;
                return this.createThumbnail(file, width, height, resizeMethod, true, function(dataUrl, canvas) {
                    if (canvas == null) // The image has not been resized
                    return callback(file);
                    else {
                        var resizeMimeType = _this.options.resizeMimeType;
                        if (resizeMimeType == null) resizeMimeType = file.type;
                        var resizedDataURL = canvas.toDataURL(resizeMimeType, _this.options.resizeQuality);
                        if (resizeMimeType === "image/jpeg" || resizeMimeType === "image/jpg") // Now add the original EXIF information
                        resizedDataURL = restoreExif(file.dataURL, resizedDataURL);
                        return callback(Dropzone.dataURItoBlob(resizedDataURL));
                    }
                }, true);
            }
        },
        {
            key: "createThumbnail",
            value: function createThumbnail(file, width, height, resizeMethod, fixOrientation, callback) {
                var _this = this;
                var ignoreExif = arguments.length > 6 && arguments[6] !== void 0 ? arguments[6] : false;
                var fileReader = new FileReader();
                fileReader.onload = function() {
                    file.dataURL = fileReader.result;
                    // Don't bother creating a thumbnail for SVG images since they're vector
                    if (file.type === "image/svg+xml") {
                        if (callback != null) callback(fileReader.result);
                        return;
                    }
                    _this.createThumbnailFromUrl(file, width, height, resizeMethod, fixOrientation, callback, undefined, ignoreExif);
                };
                fileReader.readAsDataURL(file);
            }
        },
        {
            // `mockFile` needs to have these attributes:
            //
            //     { name: 'name', size: 12345, imageUrl: '' }
            //
            // `callback` will be invoked when the image has been downloaded and displayed.
            // `crossOrigin` will be added to the `img` tag when accessing the file.
            key: "displayExistingFile",
            value: function displayExistingFile(mockFile, imageUrl, callback, crossOrigin) {
                var _this = this;
                var resizeThumbnail = arguments.length > 4 && arguments[4] !== void 0 ? arguments[4] : true;
                this.emit("addedfile", mockFile);
                this.emit("complete", mockFile);
                if (!resizeThumbnail) {
                    this.emit("thumbnail", mockFile, imageUrl);
                    if (callback) callback();
                } else {
                    var onDone = function(thumbnail) {
                        _this.emit("thumbnail", mockFile, thumbnail);
                        if (callback) callback();
                    };
                    mockFile.dataURL = imageUrl;
                    this.createThumbnailFromUrl(mockFile, this.options.thumbnailWidth, this.options.thumbnailHeight, this.options.thumbnailMethod, this.options.fixOrientation, onDone, crossOrigin);
                }
            }
        },
        {
            key: "createThumbnailFromUrl",
            value: function createThumbnailFromUrl(file, width, height, resizeMethod, fixOrientation, callback, crossOrigin) {
                var _this = this;
                var ignoreExif = arguments.length > 7 && arguments[7] !== void 0 ? arguments[7] : false;
                // Not using `new Image` here because of a bug in latest Chrome versions.
                // See https://github.com/enyo/dropzone/pull/226
                var img = document.createElement("img");
                if (crossOrigin) img.crossOrigin = crossOrigin;
                // fixOrientation is not needed anymore with browsers handling imageOrientation
                fixOrientation = getComputedStyle(document.body)["imageOrientation"] == "from-image" ? false : fixOrientation;
                img.onload = function() {
                    var loadExif = function(callback) {
                        return callback(1);
                    };
                    if (typeof EXIF !== "undefined" && EXIF !== null && fixOrientation) loadExif = function(callback) {
                        return EXIF.getData(img, function() {
                            return callback(EXIF.getTag(this, "Orientation"));
                        });
                    };
                    return loadExif(function(orientation) {
                        file.width = img.width;
                        file.height = img.height;
                        var resizeInfo = _this.options.resize.call(_this, file, width, height, resizeMethod);
                        var canvas = document.createElement("canvas");
                        var ctx = canvas.getContext("2d");
                        canvas.width = resizeInfo.trgWidth;
                        canvas.height = resizeInfo.trgHeight;
                        if (orientation > 4) {
                            canvas.width = resizeInfo.trgHeight;
                            canvas.height = resizeInfo.trgWidth;
                        }
                        switch(orientation){
                            case 2:
                                // horizontal flip
                                ctx.translate(canvas.width, 0);
                                ctx.scale(-1, 1);
                                break;
                            case 3:
                                // 180° rotate left
                                ctx.translate(canvas.width, canvas.height);
                                ctx.rotate(Math.PI);
                                break;
                            case 4:
                                // vertical flip
                                ctx.translate(0, canvas.height);
                                ctx.scale(1, -1);
                                break;
                            case 5:
                                // vertical flip + 90 rotate right
                                ctx.rotate(0.5 * Math.PI);
                                ctx.scale(1, -1);
                                break;
                            case 6:
                                // 90° rotate right
                                ctx.rotate(0.5 * Math.PI);
                                ctx.translate(0, -canvas.width);
                                break;
                            case 7:
                                // horizontal flip + 90 rotate right
                                ctx.rotate(0.5 * Math.PI);
                                ctx.translate(canvas.height, -canvas.width);
                                ctx.scale(-1, 1);
                                break;
                            case 8:
                                // 90° rotate left
                                ctx.rotate(-0.5 * Math.PI);
                                ctx.translate(-canvas.height, 0);
                                break;
                        }
                        // This is a bugfix for iOS' scaling bug.
                        drawImageIOSFix(ctx, img, resizeInfo.srcX != null ? resizeInfo.srcX : 0, resizeInfo.srcY != null ? resizeInfo.srcY : 0, resizeInfo.srcWidth, resizeInfo.srcHeight, resizeInfo.trgX != null ? resizeInfo.trgX : 0, resizeInfo.trgY != null ? resizeInfo.trgY : 0, resizeInfo.trgWidth, resizeInfo.trgHeight);
                        var thumbnail = canvas.toDataURL("image/png");
                        if (callback != null) return callback(thumbnail, canvas);
                    });
                };
                if (callback != null) img.onerror = callback;
                var dataURL = file.dataURL;
                if (ignoreExif) dataURL = removeExif(dataURL);
                return img.src = dataURL;
            }
        },
        {
            // Goes through the queue and processes files if there aren't too many already.
            key: "processQueue",
            value: function processQueue() {
                var parallelUploads = this.options.parallelUploads;
                var processingLength = this.getUploadingFiles().length;
                var i = processingLength;
                // There are already at least as many files uploading than should be
                if (processingLength >= parallelUploads) return;
                var queuedFiles = this.getQueuedFiles();
                if (!(queuedFiles.length > 0)) return;
                if (this.options.uploadMultiple) // The files should be uploaded in one request
                return this.processFiles(queuedFiles.slice(0, parallelUploads - processingLength));
                else while(i < parallelUploads){
                    if (!queuedFiles.length) return;
                     // Nothing left to process
                    this.processFile(queuedFiles.shift());
                    i++;
                }
            }
        },
        {
            // Wrapper for `processFiles`
            key: "processFile",
            value: function processFile(file) {
                return this.processFiles([
                    file
                ]);
            }
        },
        {
            // Loads the file, then calls finishedLoading()
            key: "processFiles",
            value: function processFiles(files) {
                var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                try {
                    for(var _iterator = files[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                        var file = _step.value;
                        file.processing = true; // Backwards compatibility
                        file.status = Dropzone.UPLOADING;
                        this.emit("processing", file);
                    }
                } catch (err) {
                    _didIteratorError = true;
                    _iteratorError = err;
                } finally{
                    try {
                        if (!_iteratorNormalCompletion && _iterator.return != null) {
                            _iterator.return();
                        }
                    } finally{
                        if (_didIteratorError) {
                            throw _iteratorError;
                        }
                    }
                }
                if (this.options.uploadMultiple) this.emit("processingmultiple", files);
                return this.uploadFiles(files);
            }
        },
        {
            key: "_getFilesWithXhr",
            value: function _getFilesWithXhr(xhr) {
                var files;
                return files = this.files.filter(function(file) {
                    return file.xhr === xhr;
                }).map(function(file) {
                    return file;
                });
            }
        },
        {
            // Cancels the file upload and sets the status to CANCELED
            // **if** the file is actually being uploaded.
            // If it's still in the queue, the file is being removed from it and the status
            // set to CANCELED.
            key: "cancelUpload",
            value: function cancelUpload(file) {
                if (file.status === Dropzone.UPLOADING) {
                    var groupedFiles = this._getFilesWithXhr(file.xhr);
                    var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                    try {
                        for(var _iterator = groupedFiles[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                            var groupedFile = _step.value;
                            groupedFile.status = Dropzone.CANCELED;
                        }
                    } catch (err) {
                        _didIteratorError = true;
                        _iteratorError = err;
                    } finally{
                        try {
                            if (!_iteratorNormalCompletion && _iterator.return != null) {
                                _iterator.return();
                            }
                        } finally{
                            if (_didIteratorError) {
                                throw _iteratorError;
                            }
                        }
                    }
                    if (typeof file.xhr !== "undefined") file.xhr.abort();
                    var _iteratorNormalCompletion1 = true, _didIteratorError1 = false, _iteratorError1 = undefined;
                    try {
                        for(var _iterator1 = groupedFiles[Symbol.iterator](), _step1; !(_iteratorNormalCompletion1 = (_step1 = _iterator1.next()).done); _iteratorNormalCompletion1 = true){
                            var groupedFile1 = _step1.value;
                            this.emit("canceled", groupedFile1);
                        }
                    } catch (err) {
                        _didIteratorError1 = true;
                        _iteratorError1 = err;
                    } finally{
                        try {
                            if (!_iteratorNormalCompletion1 && _iterator1.return != null) {
                                _iterator1.return();
                            }
                        } finally{
                            if (_didIteratorError1) {
                                throw _iteratorError1;
                            }
                        }
                    }
                    if (this.options.uploadMultiple) this.emit("canceledmultiple", groupedFiles);
                } else if (file.status === Dropzone.ADDED || file.status === Dropzone.QUEUED) {
                    file.status = Dropzone.CANCELED;
                    this.emit("canceled", file);
                    if (this.options.uploadMultiple) this.emit("canceledmultiple", [
                        file
                    ]);
                }
                if (this.options.autoProcessQueue) return this.processQueue();
            }
        },
        {
            key: "resolveOption",
            value: function resolveOption(option) {
                for(var _len = arguments.length, args = new Array(_len > 1 ? _len - 1 : 0), _key = 1; _key < _len; _key++){
                    args[_key - 1] = arguments[_key];
                }
                if (typeof option === "function") return option.apply(this, args);
                return option;
            }
        },
        {
            key: "uploadFile",
            value: function uploadFile(file) {
                return this.uploadFiles([
                    file
                ]);
            }
        },
        {
            key: "uploadFiles",
            value: function uploadFiles(files) {
                var _this = this;
                this._transformFiles(files, function(transformedFiles) {
                    if (_this.options.chunking) {
                        // Chunking is not allowed to be used with `uploadMultiple` so we know
                        // that there is only __one__file.
                        var transformedFile = transformedFiles[0];
                        files[0].upload.chunked = _this.options.chunking && (_this.options.forceChunking || transformedFile.size > _this.options.chunkSize);
                        files[0].upload.totalChunkCount = Math.ceil(transformedFile.size / _this.options.chunkSize);
                        if (transformedFile.size === 0) files[0].upload.totalChunkCount = 1;
                    }
                    if (files[0].upload.chunked) {
                        // This file should be sent in chunks!
                        // If the chunking option is set, we **know** that there can only be **one** file, since
                        // uploadMultiple is not allowed with this option.
                        var file = files[0];
                        var transformedFile1 = transformedFiles[0];
                        file.upload.chunks = [];
                        var handleNextChunk = function() {
                            var chunkIndex = 0;
                            // Find the next item in file.upload.chunks that is not defined yet.
                            while(file.upload.chunks[chunkIndex] !== undefined)chunkIndex++;
                            // This means, that all chunks have already been started.
                            if (chunkIndex >= file.upload.totalChunkCount) return;
                            var start = chunkIndex * _this.options.chunkSize;
                            var end = Math.min(start + _this.options.chunkSize, transformedFile1.size);
                            var dataBlock = {
                                name: _this._getParamName(0),
                                data: transformedFile1.webkitSlice ? transformedFile1.webkitSlice(start, end) : transformedFile1.slice(start, end),
                                filename: file.upload.filename,
                                chunkIndex: chunkIndex
                            };
                            file.upload.chunks[chunkIndex] = {
                                file: file,
                                index: chunkIndex,
                                dataBlock: dataBlock,
                                status: Dropzone.UPLOADING,
                                progress: 0,
                                retries: 0
                            };
                            _this._uploadData(files, [
                                dataBlock
                            ]);
                        };
                        file.upload.finishedChunkUpload = function(chunk, response) {
                            var allFinished = true;
                            chunk.status = Dropzone.SUCCESS;
                            // Clear the data from the chunk
                            chunk.dataBlock = null;
                            chunk.response = chunk.xhr.responseText;
                            chunk.responseHeaders = chunk.xhr.getAllResponseHeaders();
                            // Leaving this reference to xhr will cause memory leaks.
                            chunk.xhr = null;
                            for(var i = 0; i < file.upload.totalChunkCount; i++){
                                if (file.upload.chunks[i] === undefined) return handleNextChunk();
                                if (file.upload.chunks[i].status !== Dropzone.SUCCESS) allFinished = false;
                            }
                            if (allFinished) _this.options.chunksUploaded(file, function() {
                                _this._finished(files, response, null);
                            });
                        };
                        if (_this.options.parallelChunkUploads) {
                            // we want to limit parallelChunkUploads to the same value as parallelUploads option
                            var parallelCount = Math.min(_this.options.parallelChunkUploads === true ? _this.options.parallelUploads : _this.options.parallelChunkUploads, file.upload.totalChunkCount);
                            for(var i = 0; i < parallelCount; i++)handleNextChunk();
                        } else handleNextChunk();
                    } else {
                        var dataBlocks = [];
                        for(var i1 = 0; i1 < files.length; i1++)dataBlocks[i1] = {
                            name: _this._getParamName(i1),
                            data: transformedFiles[i1],
                            filename: files[i1].upload.filename
                        };
                        _this._uploadData(files, dataBlocks);
                    }
                });
            }
        },
        {
            /// Returns the right chunk for given file and xhr
            key: "_getChunk",
            value: function _getChunk(file, xhr) {
                for(var i = 0; i < file.upload.totalChunkCount; i++){
                    if (file.upload.chunks[i] !== undefined && file.upload.chunks[i].xhr === xhr) return file.upload.chunks[i];
                }
            }
        },
        {
            // This function actually uploads the file(s) to the server.
            //
            //  If dataBlocks contains the actual data to upload (meaning, that this could
            // either be transformed files, or individual chunks for chunked upload) then
            // they will be used for the actual data to upload.
            key: "_uploadData",
            value: function _uploadData(files, dataBlocks) {
                var _this = this;
                var xhr = new XMLHttpRequest();
                var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                try {
                    // Put the xhr object in the file objects to be able to reference it later.
                    for(var _iterator = files[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                        var file = _step.value;
                        file.xhr = xhr;
                    }
                } catch (err) {
                    _didIteratorError = true;
                    _iteratorError = err;
                } finally{
                    try {
                        if (!_iteratorNormalCompletion && _iterator.return != null) {
                            _iterator.return();
                        }
                    } finally{
                        if (_didIteratorError) {
                            throw _iteratorError;
                        }
                    }
                }
                if (files[0].upload.chunked) // Put the xhr object in the right chunk object, so it can be associated
                // later, and found with _getChunk.
                files[0].upload.chunks[dataBlocks[0].chunkIndex].xhr = xhr;
                var method = this.resolveOption(this.options.method, files, dataBlocks);
                var url = this.resolveOption(this.options.url, files, dataBlocks);
                xhr.open(method, url, true);
                // Setting the timeout after open because of IE11 issue: https://gitlab.com/meno/dropzone/issues/8
                var timeout = this.resolveOption(this.options.timeout, files);
                if (timeout) xhr.timeout = this.resolveOption(this.options.timeout, files);
                // Has to be after `.open()`. See https://github.com/enyo/dropzone/issues/179
                xhr.withCredentials = !!this.options.withCredentials;
                xhr.onload = function(e) {
                    _this._finishedUploading(files, xhr, e);
                };
                xhr.ontimeout = function() {
                    _this._handleUploadError(files, xhr, "Request timedout after ".concat(_this.options.timeout / 1000, " seconds"));
                };
                xhr.onerror = function() {
                    _this._handleUploadError(files, xhr);
                };
                // Some browsers do not have the .upload property
                var progressObj = xhr.upload != null ? xhr.upload : xhr;
                progressObj.onprogress = function(e) {
                    return _this._updateFilesUploadProgress(files, xhr, e);
                };
                var headers = this.options.defaultHeaders ? {
                    Accept: "application/json",
                    "Cache-Control": "no-cache",
                    "X-Requested-With": "XMLHttpRequest"
                } : {};
                if (this.options.binaryBody) headers["Content-Type"] = files[0].type;
                if (this.options.headers) Object.assign(headers, this.options.headers);
                for(var headerName in headers){
                    var headerValue = headers[headerName];
                    if (headerValue) xhr.setRequestHeader(headerName, headerValue);
                }
                if (this.options.binaryBody) {
                    var _iteratorNormalCompletion1 = true, _didIteratorError1 = false, _iteratorError1 = undefined;
                    try {
                        // Since the file is going to be sent as binary body, it doesn't make
                        // any sense to generate `FormData` for it.
                        for(var _iterator1 = files[Symbol.iterator](), _step1; !(_iteratorNormalCompletion1 = (_step1 = _iterator1.next()).done); _iteratorNormalCompletion1 = true){
                            var file1 = _step1.value;
                            this.emit("sending", file1, xhr);
                        }
                    } catch (err) {
                        _didIteratorError1 = true;
                        _iteratorError1 = err;
                    } finally{
                        try {
                            if (!_iteratorNormalCompletion1 && _iterator1.return != null) {
                                _iterator1.return();
                            }
                        } finally{
                            if (_didIteratorError1) {
                                throw _iteratorError1;
                            }
                        }
                    }
                    if (this.options.uploadMultiple) this.emit("sendingmultiple", files, xhr);
                    this.submitRequest(xhr, null, files);
                } else {
                    var formData = new FormData();
                    // Adding all @options parameters
                    if (this.options.params) {
                        var additionalParams = this.options.params;
                        if (typeof additionalParams === "function") additionalParams = additionalParams.call(this, files, xhr, files[0].upload.chunked ? this._getChunk(files[0], xhr) : null);
                        for(var key in additionalParams){
                            var value = additionalParams[key];
                            if (Array.isArray(value)) // The additional parameter contains an array,
                            // so lets iterate over it to attach each value
                            // individually.
                            for(var i = 0; i < value.length; i++)formData.append(key, value[i]);
                            else formData.append(key, value);
                        }
                    }
                    var _iteratorNormalCompletion2 = true, _didIteratorError2 = false, _iteratorError2 = undefined;
                    try {
                        // Let the user add additional data if necessary
                        for(var _iterator2 = files[Symbol.iterator](), _step2; !(_iteratorNormalCompletion2 = (_step2 = _iterator2.next()).done); _iteratorNormalCompletion2 = true){
                            var file2 = _step2.value;
                            this.emit("sending", file2, xhr, formData);
                        }
                    } catch (err) {
                        _didIteratorError2 = true;
                        _iteratorError2 = err;
                    } finally{
                        try {
                            if (!_iteratorNormalCompletion2 && _iterator2.return != null) {
                                _iterator2.return();
                            }
                        } finally{
                            if (_didIteratorError2) {
                                throw _iteratorError2;
                            }
                        }
                    }
                    if (this.options.uploadMultiple) this.emit("sendingmultiple", files, xhr, formData);
                    this._addFormElementData(formData);
                    // Finally add the files
                    // Has to be last because some servers (eg: S3) expect the file to be the last parameter
                    for(var i1 = 0; i1 < dataBlocks.length; i1++){
                        var dataBlock = dataBlocks[i1];
                        formData.append(dataBlock.name, dataBlock.data, dataBlock.filename);
                    }
                    this.submitRequest(xhr, formData, files);
                }
            }
        },
        {
            // Transforms all files with this.options.transformFile and invokes done with the transformed files when done.
            key: "_transformFiles",
            value: function _transformFiles(files, done) {
                var _this = this, _loop = function(i) {
                    _this.options.transformFile.call(_this, files[i], function(transformedFile) {
                        transformedFiles[i] = transformedFile;
                        if (++doneCounter === files.length) done(transformedFiles);
                    });
                };
                var transformedFiles = [];
                // Clumsy way of handling asynchronous calls, until I get to add a proper Future library.
                var doneCounter = 0;
                for(var i = 0; i < files.length; i++)_loop(i);
            }
        },
        {
            // Takes care of adding other input elements of the form to the AJAX request
            key: "_addFormElementData",
            value: function _addFormElementData(formData) {
                var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                // Take care of other input elements
                if (this.element.tagName === "FORM") try {
                    for(var _iterator = this.element.querySelectorAll("input, textarea, select, button")[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                        var input = _step.value;
                        var inputName = input.getAttribute("name");
                        var inputType = input.getAttribute("type");
                        if (inputType) inputType = inputType.toLowerCase();
                        // If the input doesn't have a name, we can't use it.
                        if (typeof inputName === "undefined" || inputName === null) continue;
                        if (input.tagName === "SELECT" && input.hasAttribute("multiple")) {
                            var _iteratorNormalCompletion1 = true, _didIteratorError1 = false, _iteratorError1 = undefined;
                            try {
                                // Possibly multiple values
                                for(var _iterator1 = input.options[Symbol.iterator](), _step1; !(_iteratorNormalCompletion1 = (_step1 = _iterator1.next()).done); _iteratorNormalCompletion1 = true){
                                    var option = _step1.value;
                                    if (option.selected) formData.append(inputName, option.value);
                                }
                            } catch (err) {
                                _didIteratorError1 = true;
                                _iteratorError1 = err;
                            } finally{
                                try {
                                    if (!_iteratorNormalCompletion1 && _iterator1.return != null) {
                                        _iterator1.return();
                                    }
                                } finally{
                                    if (_didIteratorError1) {
                                        throw _iteratorError1;
                                    }
                                }
                            }
                        } else if (!inputType || inputType !== "checkbox" && inputType !== "radio" || input.checked) formData.append(inputName, input.value);
                    }
                } catch (err) {
                    _didIteratorError = true;
                    _iteratorError = err;
                } finally{
                    try {
                        if (!_iteratorNormalCompletion && _iterator.return != null) {
                            _iterator.return();
                        }
                    } finally{
                        if (_didIteratorError) {
                            throw _iteratorError;
                        }
                    }
                }
            }
        },
        {
            // Invoked when there is new progress information about given files.
            // If e is not provided, it is assumed that the upload is finished.
            key: "_updateFilesUploadProgress",
            value: function _updateFilesUploadProgress(files, xhr, e) {
                var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                if (!files[0].upload.chunked) try {
                    // Handle file uploads without chunking
                    for(var _iterator = files[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                        var file = _step.value;
                        if (file.upload.total && file.upload.bytesSent && file.upload.bytesSent == file.upload.total) continue;
                        if (e) {
                            file.upload.progress = 100 * e.loaded / e.total;
                            file.upload.total = e.total;
                            file.upload.bytesSent = e.loaded;
                        } else {
                            // No event, so we're at 100%
                            file.upload.progress = 100;
                            file.upload.bytesSent = file.upload.total;
                        }
                        this.emit("uploadprogress", file, file.upload.progress, file.upload.bytesSent);
                    }
                } catch (err) {
                    _didIteratorError = true;
                    _iteratorError = err;
                } finally{
                    try {
                        if (!_iteratorNormalCompletion && _iterator.return != null) {
                            _iterator.return();
                        }
                    } finally{
                        if (_didIteratorError) {
                            throw _iteratorError;
                        }
                    }
                }
                else {
                    // Handle chunked file uploads
                    // Chunked upload is not compatible with uploading multiple files in one
                    // request, so we know there's only one file.
                    var file1 = files[0];
                    // Since this is a chunked upload, we need to update the appropriate chunk
                    // progress.
                    var chunk = this._getChunk(file1, xhr);
                    if (e) {
                        chunk.progress = 100 * e.loaded / e.total;
                        chunk.total = e.total;
                        chunk.bytesSent = e.loaded;
                    } else {
                        // No event, so we're at 100%
                        chunk.progress = 100;
                        chunk.bytesSent = chunk.total;
                    }
                    // Now tally the *file* upload progress from its individual chunks
                    file1.upload.progress = 0;
                    file1.upload.total = 0;
                    file1.upload.bytesSent = 0;
                    for(var i = 0; i < file1.upload.totalChunkCount; i++)if (file1.upload.chunks[i] && typeof file1.upload.chunks[i].progress !== "undefined") {
                        file1.upload.progress += file1.upload.chunks[i].progress;
                        file1.upload.total += file1.upload.chunks[i].total;
                        file1.upload.bytesSent += file1.upload.chunks[i].bytesSent;
                    }
                    // Since the process is a percentage, we need to divide by the amount of
                    // chunks we've used.
                    file1.upload.progress = file1.upload.progress / file1.upload.totalChunkCount;
                    this.emit("uploadprogress", file1, file1.upload.progress, file1.upload.bytesSent);
                }
            }
        },
        {
            key: "_finishedUploading",
            value: function _finishedUploading(files, xhr, e) {
                var response;
                if (files[0].status === Dropzone.CANCELED) return;
                if (xhr.readyState !== 4) return;
                if (xhr.responseType !== "arraybuffer" && xhr.responseType !== "blob") {
                    response = xhr.responseText;
                    if (xhr.getResponseHeader("content-type") && ~xhr.getResponseHeader("content-type").indexOf("application/json")) try {
                        response = JSON.parse(response);
                    } catch (error) {
                        e = error;
                        response = "Invalid JSON response from server.";
                    }
                }
                this._updateFilesUploadProgress(files, xhr);
                if (!(200 <= xhr.status && xhr.status < 300)) this._handleUploadError(files, xhr, response);
                else if (files[0].upload.chunked) files[0].upload.finishedChunkUpload(this._getChunk(files[0], xhr), response);
                else this._finished(files, response, e);
            }
        },
        {
            key: "_handleUploadError",
            value: function _handleUploadError(files, xhr, response) {
                if (files[0].status === Dropzone.CANCELED) return;
                if (files[0].upload.chunked && this.options.retryChunks) {
                    var chunk = this._getChunk(files[0], xhr);
                    if (chunk.retries++ < this.options.retryChunksLimit) {
                        this._uploadData(files, [
                            chunk.dataBlock
                        ]);
                        return;
                    } else console.warn("Retried this chunk too often. Giving up.");
                }
                this._errorProcessing(files, response || this.options.dictResponseError.replace("{{statusCode}}", xhr.status), xhr);
            }
        },
        {
            key: "submitRequest",
            value: function submitRequest(xhr, formData, files) {
                if (xhr.readyState != 1) {
                    console.warn("Cannot send this request because the XMLHttpRequest.readyState is not OPENED.");
                    return;
                }
                this._watchXhrSend(xhr, files); // [+]
                if (this.options.binaryBody) {
                    if (files[0].upload.chunked) {
                        var chunk = this._getChunk(files[0], xhr);
                        xhr.send(chunk.dataBlock.data);
                    } else xhr.send(files[0]);
                } else xhr.send(formData);
            }
        },
        {
            // Called internally when processing is finished.
            // Individual callbacks have to be called in the appropriate sections.
            key: "_finished",
            value: function _finished(files, responseText, e) {
                var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                try {
                    for(var _iterator = files[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                        var file = _step.value;
                        file.status = Dropzone.SUCCESS;
                        this.emit("success", file, responseText, e);
                        this.emit("complete", file);
                    }
                } catch (err) {
                    _didIteratorError = true;
                    _iteratorError = err;
                } finally{
                    try {
                        if (!_iteratorNormalCompletion && _iterator.return != null) {
                            _iterator.return();
                        }
                    } finally{
                        if (_didIteratorError) {
                            throw _iteratorError;
                        }
                    }
                }
                if (this.options.uploadMultiple) {
                    this.emit("successmultiple", files, responseText, e);
                    this.emit("completemultiple", files);
                }
                if (this.options.autoProcessQueue) return this.processQueue();
            }
        },
        {
            // Called internally when processing is finished.
            // Individual callbacks have to be called in the appropriate sections.
            key: "_errorProcessing",
            value: function _errorProcessing(files, message, xhr) {
                var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                try {
                    for(var _iterator = files[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                        var file = _step.value;
                        file.status = Dropzone.ERROR;
                        this.emit("error", file, message, xhr);
                        this.emit("complete", file);
                    }
                } catch (err) {
                    _didIteratorError = true;
                    _iteratorError = err;
                } finally{
                    try {
                        if (!_iteratorNormalCompletion && _iterator.return != null) {
                            _iterator.return();
                        }
                    } finally{
                        if (_didIteratorError) {
                            throw _iteratorError;
                        }
                    }
                }
                if (this.options.uploadMultiple) {
                    this.emit("errormultiple", files, message, xhr);
                    this.emit("completemultiple", files);
                }
                if (this.options.autoProcessQueue) return this.processQueue();
            }
        },
        {
            // ----------------------------------- //
            // [+] Custom functions: coscms.com    //
            // ----------------------------------- //
            key: "_getBytesSent",
            value: function _getBytesSent(files) {
                if (files[0].upload.chunked) return files[0].upload.bytesSent;
                var bytesSent = 0;
                var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                try {
                    for(var _iterator = files[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                        var file = _step.value;
                        bytesSent += file.upload.bytesSent;
                    }
                } catch (err) {
                    _didIteratorError = true;
                    _iteratorError = err;
                } finally{
                    try {
                        if (!_iteratorNormalCompletion && _iterator.return != null) {
                            _iterator.return();
                        }
                    } finally{
                        if (_didIteratorError) {
                            throw _iteratorError;
                        }
                    }
                }
                return bytesSent;
            }
        },
        {
            key: "_watchXhrSend",
            value: function _watchXhrSend(xhr, files, timeout) {
                var _this = this;
                if (timeout === null) timeout = 2000;
                if (timeout <= 0) return;
                var bytesSent = this._getBytesSent(files);
                var interval = setInterval(function() {
                    var bytesSentNow = _this._getBytesSent(files);
                    if (bytesSentNow !== bytesSent) {
                        bytesSent = bytesSentNow;
                        return;
                    }
                    clearInterval(interval);
                    xhr.onerror = null;
                    xhr.abort();
                    _this._handleUploadError(files, xhr, "Request timedout after ".concat(timeout, " seconds"));
                }, timeout);
                xhr.onreadystatechange = function() {
                    if (xhr.readyState == 4 || xhr.readyState == 0) try {
                        clearInterval(interval);
                    } catch (e) {}
                };
            }
        }
    ], [
        {
            key: "initClass",
            value: function initClass() {
                // Exposing the emitter class, mainly for tests
                this.prototype.Emitter = (0, _emitterDefault.default);
                /*
     This is a list of all available events you can register on a dropzone object.

     You can register an event handler like this:

     dropzone.on("dragEnter", function() { });

     */ this.prototype.events = [
                    "drop",
                    "dragstart",
                    "dragend",
                    "dragenter",
                    "dragover",
                    "dragleave",
                    "addedfile",
                    "addedfiles",
                    "removedfile",
                    "thumbnail",
                    "error",
                    "errormultiple",
                    "processing",
                    "processingmultiple",
                    "uploadprogress",
                    "totaluploadprogress",
                    "sending",
                    "sendingmultiple",
                    "success",
                    "successmultiple",
                    "canceled",
                    "canceledmultiple",
                    "complete",
                    "completemultiple",
                    "reset",
                    "maxfilesexceeded",
                    "maxfilesreached",
                    "queuecomplete"
                ];
                this.prototype._thumbnailQueue = [];
                this.prototype._processingThumbnail = false;
            }
        },
        {
            key: "uuidv4",
            value: function uuidv4() {
                return "10000000-1000-4000-8000-100000000000".replace(/[018]/g, function(c) {
                    return (+c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> +c / 4).toString(16);
                });
            }
        }
    ]);
    return Dropzone;
}((0, _emitterDefault.default));
Dropzone.initClass();
// This is a map of options for your different dropzones. Add configurations
// to this object for your different dropzone elements.
//
// Example:
//
//     Dropzone.options.myDropzoneElementId = { maxFilesize: 1 };
//
// And in html:
//
//     <form action="/upload" id="my-dropzone-element-id" class="dropzone"></form>
Dropzone.options = {};
// Returns the options for an element or undefined if none available.
Dropzone.optionsForElement = function(element) {
    // Get the `Dropzone.options.elementId` for this element if it exists
    if (element.getAttribute("id") && typeof Dropzone.options !== "undefined") return Dropzone.options[camelize(element.getAttribute("id"))];
    else return undefined;
};
// Holds a list of all dropzone instances
Dropzone.instances = [];
// Returns the dropzone for given element if any
Dropzone.forElement = function(element) {
    if (typeof element === "string") element = document.querySelector(element);
    if ((element != null ? element.dropzone : undefined) == null) throw new Error("No Dropzone found for given element. This is probably because you're trying to access it before Dropzone had the time to initialize. Use the `init` option to setup any additional observers on your Dropzone.");
    return element.dropzone;
};
// Looks for all .dropzone elements and creates a dropzone for them
Dropzone.discover = function() {
    var dropzones;
    if (document.querySelectorAll) dropzones = document.querySelectorAll(".dropzone");
    else {
        dropzones = [];
        // IE :(
        var checkElements = function(elements) {
            return function() {
                var result = [];
                var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                try {
                    for(var _iterator = elements[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                        var el = _step.value;
                        if (/(^| )dropzone($| )/.test(el.className)) result.push(dropzones.push(el));
                        else result.push(undefined);
                    }
                } catch (err) {
                    _didIteratorError = true;
                    _iteratorError = err;
                } finally{
                    try {
                        if (!_iteratorNormalCompletion && _iterator.return != null) {
                            _iterator.return();
                        }
                    } finally{
                        if (_didIteratorError) {
                            throw _iteratorError;
                        }
                    }
                }
                return result;
            }();
        };
        checkElements(document.getElementsByTagName("div"));
        checkElements(document.getElementsByTagName("form"));
    }
    return function() {
        var result = [];
        var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
        try {
            for(var _iterator = dropzones[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                var dropzone = _step.value;
                // Create a dropzone unless auto discover has been disabled for specific element
                if (Dropzone.optionsForElement(dropzone) !== false) result.push(new Dropzone(dropzone));
                else result.push(undefined);
            }
        } catch (err) {
            _didIteratorError = true;
            _iteratorError = err;
        } finally{
            try {
                if (!_iteratorNormalCompletion && _iterator.return != null) {
                    _iterator.return();
                }
            } finally{
                if (_didIteratorError) {
                    throw _iteratorError;
                }
            }
        }
        return result;
    }();
};
// Checks if the browser is supported by simply checking if Promise is here: a good cutoff
Dropzone.isBrowserSupported = function() {
    return typeof Promise !== "undefined";
};
Dropzone.dataURItoBlob = function(dataURI) {
    // convert base64 to raw binary data held in a string
    // doesn't handle URLEncoded DataURIs - see SO answer #6850276 for code that does this
    var byteString = atob(dataURI.split(",")[1]);
    // separate out the mime component
    var mimeString = dataURI.split(",")[0].split(":")[1].split(";")[0];
    // write the bytes of the string to an ArrayBuffer
    var ab = new ArrayBuffer(byteString.length);
    var ia = new Uint8Array(ab);
    for(var i = 0, end = byteString.length, asc = 0 <= end; asc ? i <= end : i >= end; asc ? i++ : i--)ia[i] = byteString.charCodeAt(i);
    // write the ArrayBuffer to a blob
    return new Blob([
        ab
    ], {
        type: mimeString
    });
};
// Returns an array without the rejected item
var without = function(list, rejectedItem) {
    return list.filter(function(item) {
        return item !== rejectedItem;
    }).map(function(item) {
        return item;
    });
};
// abc-def_ghi -> abcDefGhi
var camelize = function(str) {
    return str.replace(/[\-_](\w)/g, function(match) {
        return match.charAt(1).toUpperCase();
    });
};
// Creates an element from string
Dropzone.createElement = function(string) {
    var div = document.createElement("div");
    div.innerHTML = string;
    return div.childNodes[0];
};
// Tests if given element is inside (or simply is) the container
Dropzone.elementInside = function(element, container) {
    if (element === container) return true;
     // Coffeescript doesn't support do/while loops
    while(element = element.parentNode){
        if (element === container) return true;
    }
    return false;
};
Dropzone.getElement = function(el, name) {
    var element;
    if (typeof el === "string") element = document.querySelector(el);
    else if (el.nodeType != null) element = el;
    if (element == null) throw new Error("Invalid `".concat(name, "` option provided. Please provide a CSS selector or a plain HTML element."));
    return element;
};
Dropzone.getElements = function(els, name) {
    var el, elements;
    if (els instanceof Array) {
        elements = [];
        try {
            var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
            try {
                for(var _iterator = els[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                    el = _step.value;
                    elements.push(this.getElement(el, name));
                }
            } catch (err) {
                _didIteratorError = true;
                _iteratorError = err;
            } finally{
                try {
                    if (!_iteratorNormalCompletion && _iterator.return != null) {
                        _iterator.return();
                    }
                } finally{
                    if (_didIteratorError) {
                        throw _iteratorError;
                    }
                }
            }
        } catch (e) {
            elements = null;
        }
    } else if (typeof els === "string") {
        elements = [];
        var _iteratorNormalCompletion1 = true, _didIteratorError1 = false, _iteratorError1 = undefined;
        try {
            for(var _iterator1 = document.querySelectorAll(els)[Symbol.iterator](), _step1; !(_iteratorNormalCompletion1 = (_step1 = _iterator1.next()).done); _iteratorNormalCompletion1 = true){
                el = _step1.value;
                elements.push(el);
            }
        } catch (err) {
            _didIteratorError1 = true;
            _iteratorError1 = err;
        } finally{
            try {
                if (!_iteratorNormalCompletion1 && _iterator1.return != null) {
                    _iterator1.return();
                }
            } finally{
                if (_didIteratorError1) {
                    throw _iteratorError1;
                }
            }
        }
    } else if (els.nodeType != null) elements = [
        els
    ];
    if (elements == null || !elements.length) throw new Error("Invalid `".concat(name, "` option provided. Please provide a CSS selector, a plain HTML element or a list of those."));
    return elements;
};
// Asks the user the question and calls accepted or rejected accordingly
//
// The default implementation just uses `window.confirm` and then calls the
// appropriate callback.
Dropzone.confirm = function(question, accepted, rejected) {
    if (window.confirm(question)) return accepted();
    else if (rejected != null) return rejected();
};
// Validates the mime type like this:
//
// https://developer.mozilla.org/en-US/docs/HTML/Element/input#attr-accept
Dropzone.isValidFile = function(file, acceptedFiles) {
    if (!acceptedFiles) return true;
     // If there are no accepted mime types, it's OK
    acceptedFiles = acceptedFiles.split(",");
    var mimeType = file.type;
    var baseMimeType = mimeType.replace(/\/.*$/, "");
    var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
    try {
        for(var _iterator = acceptedFiles[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
            var validType = _step.value;
            validType = validType.trim();
            if (validType.charAt(0) === ".") {
                if (file.name.toLowerCase().indexOf(validType.toLowerCase(), file.name.length - validType.length) !== -1) return true;
            } else if (/\/\*$/.test(validType)) {
                // This is something like a image/* mime type
                if (baseMimeType === validType.replace(/\/.*$/, "")) return true;
            } else {
                if (mimeType === validType) return true;
            }
        }
    } catch (err) {
        _didIteratorError = true;
        _iteratorError = err;
    } finally{
        try {
            if (!_iteratorNormalCompletion && _iterator.return != null) {
                _iterator.return();
            }
        } finally{
            if (_didIteratorError) {
                throw _iteratorError;
            }
        }
    }
    return false;
};
// Augment jQuery
if (typeof jQuery !== "undefined" && jQuery !== null) jQuery.fn.dropzone = function(options) {
    return this.each(function() {
        return new Dropzone(this, options);
    });
};
// Dropzone file status codes
Dropzone.ADDED = "added";
Dropzone.QUEUED = "queued";
// For backwards compatibility. Now, if a file is accepted, it's either queued
// or uploading.
Dropzone.ACCEPTED = Dropzone.QUEUED;
Dropzone.UPLOADING = "uploading";
Dropzone.PROCESSING = Dropzone.UPLOADING; // alias
Dropzone.CANCELED = "canceled";
Dropzone.ERROR = "error";
Dropzone.SUCCESS = "success";
/*

 Bugfix for iOS 6 and 7
 Source: http://stackoverflow.com/questions/11929099/html5-canvas-drawimage-ratio-bug-ios
 based on the work of https://github.com/stomita/ios-imagefile-megapixel

 */ // Detecting vertical squash in loaded image.
// Fixes a bug which squash image vertically while drawing into canvas for some images.
// This is a bug in iOS6 devices. This function from https://github.com/stomita/ios-imagefile-megapixel
var detectVerticalSquash = function detectVerticalSquash(img) {
    var ih = img.naturalHeight;
    var canvas = document.createElement("canvas");
    canvas.width = 1;
    canvas.height = ih;
    var ctx = canvas.getContext("2d");
    ctx.drawImage(img, 0, 0);
    var data = ctx.getImageData(1, 0, 1, ih).data;
    // search image edge pixel position in case it is squashed vertically.
    var sy = 0;
    var ey = ih;
    var py = ih;
    while(py > sy){
        var alpha = data[(py - 1) * 4 + 3];
        if (alpha === 0) ey = py;
        else sy = py;
        py = ey + sy >> 1;
    }
    var ratio = py / ih;
    if (ratio === 0) return 1;
    else return ratio;
};
// A replacement for context.drawImage
// (args are for source and destination).
var drawImageIOSFix = function drawImageIOSFix(ctx, img, sx, sy, sw, sh, dx, dy, dw, dh) {
    var vertSquashRatio = detectVerticalSquash(img);
    return ctx.drawImage(img, sx, sy, sw, sh, dx, dy, dw, dh / vertSquashRatio);
};
// Inspired by MinifyJpeg
// Source: http://www.perry.cz/files/ExifRestorer.js
// http://elicon.blog57.fc2.com/blog-entry-206.html
function removeExif(origFileBase64) {
    var marker = "data:image/jpeg;base64,";
    if (!origFileBase64.startsWith(marker)) return origFileBase64;
    var origFile = window.atob(origFileBase64.slice(marker.length));
    if (!origFile.startsWith("\xff\xd8\xff")) return origFileBase64;
    // loop through the JPEG file segments and copy all but Exif segments into the filtered file.
    var head = 0;
    var filteredFile = "";
    while(head < origFile.length){
        if (origFile.slice(head, head + 2) == "\xff\xda") {
            // this is the start of the image data, we don't expect exif data after that.
            filteredFile += origFile.slice(head);
            break;
        } else if (origFile.slice(head, head + 2) == "\xff\xd8") {
            // this is the global start marker.
            filteredFile += origFile.slice(head, head + 2);
            head += 2;
        } else {
            // we have a segment of variable size.
            var length = origFile.charCodeAt(head + 2) * 256 + origFile.charCodeAt(head + 3);
            var endPoint = head + length + 2;
            var segment = origFile.slice(head, endPoint);
            if (!segment.startsWith("\xff\xe1")) filteredFile += segment;
            head = endPoint;
        }
    }
    return marker + window.btoa(filteredFile);
}
function restoreExif(origFileBase64, resizedFileBase64) {
    var marker = "data:image/jpeg;base64,";
    if (!(origFileBase64.startsWith(marker) && resizedFileBase64.startsWith(marker))) return resizedFileBase64;
    var origFile = window.atob(origFileBase64.slice(marker.length));
    if (!origFile.startsWith("\xff\xd8\xff")) return resizedFileBase64;
    // Go through the JPEG file segments one by one and collect any Exif segments we find.
    var head = 0;
    var exifData = "";
    while(head < origFile.length){
        if (origFile.slice(head, head + 2) == "\xff\xda") break;
        else if (origFile.slice(head, head + 2) == "\xff\xd8") // this is the global start marker.
        head += 2;
        else {
            // we have a segment of variable size.
            var length = origFile.charCodeAt(head + 2) * 256 + origFile.charCodeAt(head + 3);
            var endPoint = head + length + 2;
            var segment = origFile.slice(head, endPoint);
            if (segment.startsWith("\xff\xe1")) exifData += segment;
            head = endPoint;
        }
    }
    if (exifData == "") return resizedFileBase64;
    var resizedFile = window.atob(resizedFileBase64.slice(marker.length));
    if (!resizedFile.startsWith("\xff\xd8\xff")) return resizedFileBase64;
    // The first file segment is always header information so insert the Exif data as second segment.
    var splitPoint = 4 + resizedFile.charCodeAt(4) * 256 + resizedFile.charCodeAt(5);
    resizedFile = resizedFile.slice(0, splitPoint) + exifData + resizedFile.slice(splitPoint);
    return marker + window.btoa(resizedFile);
}
function __guard__(value, transform) {
    return typeof value !== "undefined" && value !== null ? transform(value) : undefined;
}
function __guardMethod__(obj, methodName, transform) {
    if (typeof obj !== "undefined" && obj !== null && typeof obj[methodName] === "function") return transform(obj, methodName);
    else return undefined;
}

},{"@swc/helpers/_/_assert_this_initialized":"4Dqij","@swc/helpers/_/_class_call_check":"INQRr","@swc/helpers/_/_create_class":"eJoMn","@swc/helpers/_/_inherits":"bEQbf","@swc/helpers/_/_object_spread":"lU0hD","@swc/helpers/_/_object_spread_props":"bWXxc","@swc/helpers/_/_possible_constructor_return":"9YtSZ","@swc/helpers/_/_create_super":"1xeF0","./emitter":"3dfRA","./options":"1Xeqc","@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"4Dqij":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_assert_this_initialized", function() {
    return _assert_this_initialized;
});
parcelHelpers.export(exports, "_", function() {
    return _assert_this_initialized;
});
function _assert_this_initialized(self) {
    if (self === void 0) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
    return self;
}

},{"@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"INQRr":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_class_call_check", function() {
    return _class_call_check;
});
parcelHelpers.export(exports, "_", function() {
    return _class_call_check;
});
function _class_call_check(instance, Constructor) {
    if (!(instance instanceof Constructor)) throw new TypeError("Cannot call a class as a function");
}

},{"@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"eJoMn":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_create_class", function() {
    return _create_class;
});
parcelHelpers.export(exports, "_", function() {
    return _create_class;
});
function _defineProperties(target, props) {
    for(var i = 0; i < props.length; i++){
        var descriptor = props[i];
        descriptor.enumerable = descriptor.enumerable || false;
        descriptor.configurable = true;
        if ("value" in descriptor) descriptor.writable = true;
        Object.defineProperty(target, descriptor.key, descriptor);
    }
}
function _create_class(Constructor, protoProps, staticProps) {
    if (protoProps) _defineProperties(Constructor.prototype, protoProps);
    if (staticProps) _defineProperties(Constructor, staticProps);
    return Constructor;
}

},{"@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"bEQbf":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_inherits", function() {
    return _inherits;
});
parcelHelpers.export(exports, "_", function() {
    return _inherits;
});
var _setPrototypeOfJs = require("./_set_prototype_of.js");
function _inherits(subClass, superClass) {
    if (typeof superClass !== "function" && superClass !== null) throw new TypeError("Super expression must either be null or a function");
    subClass.prototype = Object.create(superClass && superClass.prototype, {
        constructor: {
            value: subClass,
            writable: true,
            configurable: true
        }
    });
    if (superClass) (0, _setPrototypeOfJs._set_prototype_of)(subClass, superClass);
}

},{"./_set_prototype_of.js":"krDh8","@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"krDh8":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_set_prototype_of", function() {
    return _set_prototype_of;
});
parcelHelpers.export(exports, "_", function() {
    return _set_prototype_of;
});
function _set_prototype_of(o, p) {
    _set_prototype_of = Object.setPrototypeOf || function setPrototypeOf(o, p) {
        o.__proto__ = p;
        return o;
    };
    return _set_prototype_of(o, p);
}

},{"@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"lU0hD":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_object_spread", function() {
    return _object_spread;
});
parcelHelpers.export(exports, "_", function() {
    return _object_spread;
});
var _definePropertyJs = require("./_define_property.js");
function _object_spread(target) {
    for(var i = 1; i < arguments.length; i++){
        var source = arguments[i] != null ? arguments[i] : {};
        var ownKeys = Object.keys(source);
        if (typeof Object.getOwnPropertySymbols === "function") ownKeys = ownKeys.concat(Object.getOwnPropertySymbols(source).filter(function(sym) {
            return Object.getOwnPropertyDescriptor(source, sym).enumerable;
        }));
        ownKeys.forEach(function(key) {
            (0, _definePropertyJs._define_property)(target, key, source[key]);
        });
    }
    return target;
}

},{"./_define_property.js":"8QkZB","@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"8QkZB":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_define_property", function() {
    return _define_property;
});
parcelHelpers.export(exports, "_", function() {
    return _define_property;
});
function _define_property(obj, key, value) {
    if (key in obj) Object.defineProperty(obj, key, {
        value: value,
        enumerable: true,
        configurable: true,
        writable: true
    });
    else obj[key] = value;
    return obj;
}

},{"@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"bWXxc":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_object_spread_props", function() {
    return _object_spread_props;
});
parcelHelpers.export(exports, "_", function() {
    return _object_spread_props;
});
function ownKeys(object, enumerableOnly) {
    var keys = Object.keys(object);
    if (Object.getOwnPropertySymbols) {
        var symbols = Object.getOwnPropertySymbols(object);
        if (enumerableOnly) symbols = symbols.filter(function(sym) {
            return Object.getOwnPropertyDescriptor(object, sym).enumerable;
        });
        keys.push.apply(keys, symbols);
    }
    return keys;
}
function _object_spread_props(target, source) {
    source = source != null ? source : {};
    if (Object.getOwnPropertyDescriptors) Object.defineProperties(target, Object.getOwnPropertyDescriptors(source));
    else ownKeys(Object(source)).forEach(function(key) {
        Object.defineProperty(target, key, Object.getOwnPropertyDescriptor(source, key));
    });
    return target;
}

},{"@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"9YtSZ":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_possible_constructor_return", function() {
    return _possible_constructor_return;
});
parcelHelpers.export(exports, "_", function() {
    return _possible_constructor_return;
});
var _assertThisInitializedJs = require("./_assert_this_initialized.js");
var _typeOfJs = require("./_type_of.js");
function _possible_constructor_return(self, call) {
    if (call && ((0, _typeOfJs._type_of)(call) === "object" || typeof call === "function")) return call;
    return (0, _assertThisInitializedJs._assert_this_initialized)(self);
}

},{"./_assert_this_initialized.js":"4Dqij","./_type_of.js":"lQLoi","@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"1xeF0":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_create_super", function() {
    return _create_super;
});
parcelHelpers.export(exports, "_", function() {
    return _create_super;
});
var _getPrototypeOfJs = require("./_get_prototype_of.js");
var _isNativeReflectConstructJs = require("./_is_native_reflect_construct.js");
var _possibleConstructorReturnJs = require("./_possible_constructor_return.js");
function _create_super(Derived) {
    var hasNativeReflectConstruct = (0, _isNativeReflectConstructJs._is_native_reflect_construct)();
    return function _createSuperInternal() {
        var Super = (0, _getPrototypeOfJs._get_prototype_of)(Derived), result;
        if (hasNativeReflectConstruct) {
            var NewTarget = (0, _getPrototypeOfJs._get_prototype_of)(this).constructor;
            result = Reflect.construct(Super, arguments, NewTarget);
        } else result = Super.apply(this, arguments);
        return (0, _possibleConstructorReturnJs._possible_constructor_return)(this, result);
    };
}

},{"./_get_prototype_of.js":"hXyLP","./_is_native_reflect_construct.js":"6vhqH","./_possible_constructor_return.js":"9YtSZ","@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"hXyLP":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_get_prototype_of", function() {
    return _get_prototype_of;
});
parcelHelpers.export(exports, "_", function() {
    return _get_prototype_of;
});
function _get_prototype_of(o) {
    _get_prototype_of = Object.setPrototypeOf ? Object.getPrototypeOf : function getPrototypeOf(o) {
        return o.__proto__ || Object.getPrototypeOf(o);
    };
    return _get_prototype_of(o);
}

},{"@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"6vhqH":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "_is_native_reflect_construct", function() {
    return _is_native_reflect_construct;
});
parcelHelpers.export(exports, "_", function() {
    return _is_native_reflect_construct;
});
function _is_native_reflect_construct() {
    if (typeof Reflect === "undefined" || !Reflect.construct) return false;
    if (Reflect.construct.sham) return false;
    if (typeof Proxy === "function") return true;
    try {
        Boolean.prototype.valueOf.call(Reflect.construct(Boolean, [], function() {}));
        return true;
    } catch (e) {
        return false;
    }
}

},{"@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"3dfRA":[function(require,module,exports) {
// The Emitter class provides the ability to call `.on()` on Dropzone to listen
// to events.
// It is strongly based on component's emitter class, and I removed the
// functionality because of the dependency hell with different frameworks.
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
parcelHelpers.export(exports, "default", function() {
    return Emitter;
});
var _classCallCheck = require("@swc/helpers/_/_class_call_check");
var _createClass = require("@swc/helpers/_/_create_class");
var Emitter = /*#__PURE__*/ function() {
    "use strict";
    function Emitter() {
        (0, _classCallCheck._)(this, Emitter);
    }
    (0, _createClass._)(Emitter, [
        {
            // Add an event listener for given event
            key: "on",
            value: function on(event, fn) {
                this._callbacks = this._callbacks || {};
                // Create namespace for this event
                if (!this._callbacks[event]) this._callbacks[event] = [];
                this._callbacks[event].push(fn);
                return this;
            }
        },
        {
            key: "emit",
            value: function emit(event) {
                for(var _len = arguments.length, args = new Array(_len > 1 ? _len - 1 : 0), _key = 1; _key < _len; _key++){
                    args[_key - 1] = arguments[_key];
                }
                this._callbacks = this._callbacks || {};
                var callbacks = this._callbacks[event];
                var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
                if (callbacks) try {
                    for(var _iterator = callbacks[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                        var callback = _step.value;
                        callback.apply(this, args);
                    }
                } catch (err) {
                    _didIteratorError = true;
                    _iteratorError = err;
                } finally{
                    try {
                        if (!_iteratorNormalCompletion && _iterator.return != null) {
                            _iterator.return();
                        }
                    } finally{
                        if (_didIteratorError) {
                            throw _iteratorError;
                        }
                    }
                }
                // trigger a corresponding DOM event
                if (this.element) this.element.dispatchEvent(this.makeEvent("dropzone:" + event, {
                    args: args
                }));
                return this;
            }
        },
        {
            key: "makeEvent",
            value: function makeEvent(eventName, detail) {
                var params = {
                    bubbles: true,
                    cancelable: true,
                    detail: detail
                };
                if (typeof window.CustomEvent === "function") return new CustomEvent(eventName, params);
                else {
                    // IE 11 support
                    // https://developer.mozilla.org/en-US/docs/Web/API/CustomEvent/CustomEvent
                    var evt = document.createEvent("CustomEvent");
                    evt.initCustomEvent(eventName, params.bubbles, params.cancelable, params.detail);
                    return evt;
                }
            }
        },
        {
            // Remove event listener for given event. If fn is not provided, all event
            // listeners for that event will be removed. If neither is provided, all
            // event listeners will be removed.
            key: "off",
            value: function off(event, fn) {
                if (!this._callbacks || arguments.length === 0) {
                    this._callbacks = {};
                    return this;
                }
                // specific event
                var callbacks = this._callbacks[event];
                if (!callbacks) return this;
                // remove all handlers
                if (arguments.length === 1) {
                    delete this._callbacks[event];
                    return this;
                }
                // remove specific handler
                for(var i = 0; i < callbacks.length; i++){
                    var callback = callbacks[i];
                    if (callback === fn) {
                        callbacks.splice(i, 1);
                        break;
                    }
                }
                return this;
            }
        }
    ]);
    return Emitter;
}();

},{"@swc/helpers/_/_class_call_check":"INQRr","@swc/helpers/_/_create_class":"eJoMn","@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"1Xeqc":[function(require,module,exports) {
var parcelHelpers = require("@parcel/transformer-js/src/esmodule-helpers.js");
parcelHelpers.defineInteropFlag(exports);
var _dropzone = require("./dropzone");
var _dropzoneDefault = parcelHelpers.interopDefault(_dropzone);
var _previewTemplateHtml = require("bundle-text:./preview-template.html");
var _previewTemplateHtmlDefault = parcelHelpers.interopDefault(_previewTemplateHtml);
var defaultOptions = {
    /**
   * Has to be specified on elements other than form (or when the form doesn't
   * have an `action` attribute).
   *
   * You can also provide a function that will be called with `files` and
   * `dataBlocks`  and must return the url as string.
   */ url: null,
    /**
   * Can be changed to `"put"` if necessary. You can also provide a function
   * that will be called with `files` and must return the method (since `v3.12.0`).
   */ method: "post",
    /**
   * Will be set on the XHRequest.
   */ withCredentials: false,
    /**
   * The timeout for the XHR requests in milliseconds (since `v4.4.0`).
   * If set to null or 0, no timeout is going to be set.
   */ timeout: null,
    /**
   * How many file uploads to process in parallel (See the
   * Enqueuing file uploads documentation section for more info)
   */ parallelUploads: 2,
    /**
   * Whether to send multiple files in one request. If
   * this it set to true, then the fallback file input element will
   * have the `multiple` attribute as well. This option will
   * also trigger additional events (like `processingmultiple`). See the events
   * documentation section for more information.
   */ uploadMultiple: false,
    /**
   * Whether you want files to be uploaded in chunks to your server. This can't be
   * used in combination with `uploadMultiple`.
   *
   * See [chunksUploaded](#config-chunksUploaded) for the callback to finalise an upload.
   */ chunking: false,
    /**
   * If `chunking` is enabled, this defines whether **every** file should be chunked,
   * even if the file size is below chunkSize. This means, that the additional chunk
   * form data will be submitted and the `chunksUploaded` callback will be invoked.
   */ forceChunking: false,
    /**
   * If `chunking` is `true`, then this defines the chunk size in bytes.
   */ chunkSize: 2097152,
    /**
   * If `true`, the individual chunks of a file are being uploaded simultaneously.
   * The limit of concurrent connections is governed by `parallelUploads`.
   */ parallelChunkUploads: false,
    /**
   * Whether a chunk should be retried if it fails.
   */ retryChunks: false,
    /**
   * If `retryChunks` is true, how many times should it be retried.
   */ retryChunksLimit: 3,
    /**
   * The maximum filesize (in MiB) that is allowed to be uploaded.
   */ maxFilesize: 256,
    /**
   * The name of the file param that gets transferred.
   * **NOTE**: If you have the option  `uploadMultiple` set to `true`, then
   * Dropzone will append `[]` to the name.
   */ paramName: "file",
    /**
   * Whether thumbnails for images should be generated
   */ createImageThumbnails: true,
    /**
   * In MB. When the filename exceeds this limit, the thumbnail will not be generated.
   */ maxThumbnailFilesize: 10,
    /**
   * If `null`, the ratio of the image will be used to calculate it.
   */ thumbnailWidth: 120,
    /**
   * The same as `thumbnailWidth`. If both are null, images will not be resized.
   */ thumbnailHeight: 120,
    /**
   * How the images should be scaled down in case both, `thumbnailWidth` and `thumbnailHeight` are provided.
   * Can be either `contain` or `crop`.
   */ thumbnailMethod: "crop",
    /**
   * If set, images will be resized to these dimensions before being **uploaded**.
   * If only one, `resizeWidth` **or** `resizeHeight` is provided, the original aspect
   * ratio of the file will be preserved.
   *
   * The `options.transformFile` function uses these options, so if the `transformFile` function
   * is overridden, these options don't do anything.
   */ resizeWidth: null,
    /**
   * See `resizeWidth`.
   */ resizeHeight: null,
    /**
   * The mime type of the resized image (before it gets uploaded to the server).
   * If `null` the original mime type will be used. To force jpeg, for example, use `image/jpeg`.
   * See `resizeWidth` for more information.
   */ resizeMimeType: null,
    /**
   * The quality of the resized images. See `resizeWidth`.
   */ resizeQuality: 0.8,
    /**
   * How the images should be scaled down in case both, `resizeWidth` and `resizeHeight` are provided.
   * Can be either `contain` or `crop`.
   */ resizeMethod: "contain",
    /**
   * The base that is used to calculate the **displayed** filesize. You can
   * change this to 1024 if you would rather display kibibytes, mebibytes,
   * etc... 1024 is technically incorrect, because `1024 bytes` are `1 kibibyte`
   * not `1 kilobyte`. You can change this to `1024` if you don't care about
   * validity.
   */ filesizeBase: 1000,
    /**
   * If not `null` defines how many files this Dropzone handles. If it exceeds,
   * the event `maxfilesexceeded` will be called. The dropzone element gets the
   * class `dz-max-files-reached` accordingly so you can provide visual
   * feedback.
   */ maxFiles: null,
    /**
   * An optional object to send additional headers to the server. Eg:
   * `{ "My-Awesome-Header": "header value" }`
   */ headers: null,
    /**
   * Should the default headers be set or not?
   * Accept: application/json <- for requesting json response
   * Cache-Control: no-cache <- Request shouldn't be cached
   * X-Requested-With: XMLHttpRequest <- We sent the request via XMLHttpRequest
   */ defaultHeaders: true,
    /**
   * If `true`, the dropzone element itself will be clickable, if `false`
   * nothing will be clickable.
   *
   * You can also pass an HTML element, a CSS selector (for multiple elements)
   * or an array of those. In that case, all of those elements will trigger an
   * upload when clicked.
   */ clickable: true,
    /**
   * Whether hidden files in directories should be ignored.
   */ ignoreHiddenFiles: true,
    /**
   * The default implementation of `accept` checks the file's mime type or
   * extension against this list. This is a comma separated list of mime
   * types or file extensions.
   *
   * Eg.: `image/*,application/pdf,.psd`
   *
   * If the Dropzone is `clickable` this option will also be used as
   * [`accept`](https://developer.mozilla.org/en-US/docs/HTML/Element/input#attr-accept)
   * parameter on the hidden file input as well.
   */ acceptedFiles: null,
    /**
   * If false, files will be added to the queue but the queue will not be
   * processed automatically.
   * This can be useful if you need some additional user input before sending
   * files (or if you want want all files sent at once).
   * If you're ready to send the file simply call `myDropzone.processQueue()`.
   *
   * See the [enqueuing file uploads](#enqueuing-file-uploads) documentation
   * section for more information.
   */ autoProcessQueue: true,
    /**
   * If false, files added to the dropzone will not be queued by default.
   * You'll have to call `enqueueFile(file)` manually.
   */ autoQueue: true,
    /**
   * If `true`, this will add a link to every file preview to remove or cancel (if
   * already uploading) the file. The `dictCancelUpload`, `dictCancelUploadConfirmation`
   * and `dictRemoveFile` options are used for the wording.
   */ addRemoveLinks: false,
    /**
   * Defines where to display the file previews – if `null` the
   * Dropzone element itself is used. Can be a plain `HTMLElement` or a CSS
   * selector. The element should have the `dropzone-previews` class so
   * the previews are displayed properly.
   */ previewsContainer: null,
    /**
   * Set this to `true` if you don't want previews to be shown.
   */ disablePreviews: false,
    /**
   * This is the element the hidden input field (which is used when clicking on the
   * dropzone to trigger file selection) will be appended to. This might
   * be important in case you use frameworks to switch the content of your page.
   *
   * Can be a selector string, or an element directly.
   */ hiddenInputContainer: "body",
    /**
   * If null, no capture type will be specified
   * If camera, mobile devices will skip the file selection and choose camera
   * If microphone, mobile devices will skip the file selection and choose the microphone
   * If camcorder, mobile devices will skip the file selection and choose the camera in video mode
   * On apple devices multiple must be set to false.  AcceptedFiles may need to
   * be set to an appropriate mime type (e.g. "image/*", "audio/*", or "video/*").
   */ capture: null,
    /**
   * **Deprecated**. Use `renameFile` instead.
   */ renameFilename: null,
    /**
   * A function that is invoked before the file is uploaded to the server and renames the file.
   * This function gets the `File` as argument and can use the `file.name`. The actual name of the
   * file that gets used during the upload can be accessed through `file.upload.filename`.
   */ renameFile: null,
    /**
   * If `true` the fallback will be forced. This is very useful to test your server
   * implementations first and make sure that everything works as
   * expected without dropzone if you experience problems, and to test
   * how your fallbacks will look.
   */ forceFallback: false,
    /**
   * The text used before any files are dropped.
   */ dictDefaultMessage: "Drop files here to upload",
    /**
   * The text that replaces the default message text it the browser is not supported.
   */ dictFallbackMessage: "Your browser does not support drag'n'drop file uploads.",
    /**
   * The text that will be added before the fallback form.
   * If you provide a  fallback element yourself, or if this option is `null` this will
   * be ignored.
   */ dictFallbackText: "Please use the fallback form below to upload your files like in the olden days.",
    /**
   * If the filesize is too big.
   * `{{filesize}}` and `{{maxFilesize}}` will be replaced with the respective configuration values.
   */ dictFileTooBig: "File is too big ({{filesize}}MiB). Max filesize: {{maxFilesize}}MiB.",
    /**
   * If the file doesn't match the file type.
   */ dictInvalidFileType: "You can't upload files of this type.",
    /**
   * If the server response was invalid.
   * `{{statusCode}}` will be replaced with the servers status code.
   */ dictResponseError: "Server responded with {{statusCode}} code.",
    /**
   * If `addRemoveLinks` is true, the text to be used for the cancel upload link.
   */ dictCancelUpload: "Cancel upload",
    /**
   * The text that is displayed if an upload was manually canceled
   */ dictUploadCanceled: "Upload canceled.",
    /**
   * If `addRemoveLinks` is true, the text to be used for confirmation when cancelling upload.
   */ dictCancelUploadConfirmation: "Are you sure you want to cancel this upload?",
    /**
   * If `addRemoveLinks` is true, the text to be used to remove a file.
   */ dictRemoveFile: "Remove file",
    /**
   * If this is not null, then the user will be prompted before removing a file.
   */ dictRemoveFileConfirmation: null,
    /**
   * Displayed if `maxFiles` is st and exceeded.
   * The string `{{maxFiles}}` will be replaced by the configuration value.
   */ dictMaxFilesExceeded: "You cannot upload any more files.",
    /**
   * Allows you to translate the different units. Starting with `tb` for terabytes and going down to
   * `b` for bytes.
   */ dictFileSizeUnits: {
        tb: "TB",
        gb: "GB",
        mb: "MB",
        kb: "KB",
        b: "b"
    },
    /**
   * Called when dropzone initialized
   * You can add event listeners here
   */ init: function() {},
    /**
   * Can be an **object** of additional parameters to transfer to the server, **or** a `Function`
   * that gets invoked with the `files`, `xhr` and, if it's a chunked upload, `chunk` arguments. In case
   * of a function, this needs to return a map.
   *
   * The default implementation does nothing for normal uploads, but adds relevant information for
   * chunked uploads.
   *
   * This is the same as adding hidden input fields in the form element.
   */ params: function(files, xhr, chunk) {
        if (chunk) return {
            dzuuid: chunk.file.upload.uuid,
            dzchunkindex: chunk.index,
            dztotalfilesize: chunk.file.size,
            dzchunksize: this.options.chunkSize,
            dztotalchunkcount: chunk.file.upload.totalChunkCount,
            dzchunkbyteoffset: chunk.index * this.options.chunkSize
        };
    },
    /**
   * A function that gets a [file](https://developer.mozilla.org/en-US/docs/DOM/File)
   * and a `done` function as parameters.
   *
   * If the done function is invoked without arguments, the file is "accepted" and will
   * be processed. If you pass an error message, the file is rejected, and the error
   * message will be displayed.
   * This function will not be called if the file is too big or doesn't match the mime types.
   */ accept: function(file, done) {
        return done();
    },
    /**
   * The callback that will be invoked when all chunks have been uploaded for a file.
   * It gets the file for which the chunks have been uploaded as the first parameter,
   * and the `done` function as second. `done()` needs to be invoked when everything
   * needed to finish the upload process is done.
   */ chunksUploaded: function chunksUploaded(file, done) {
        done();
    },
    /**
   * Sends the file as binary blob in body instead of form data.
   * If this is set, the `params` option will be ignored.
   * It's an error to set this to `true` along with `uploadMultiple` since
   * multiple files cannot be in a single binary body.
   */ binaryBody: false,
    /**
   * Gets called when the browser is not supported.
   * The default implementation shows the fallback input field and adds
   * a text.
   */ fallback: function() {
        // This code should pass in IE7... :(
        var messageElement;
        this.element.className = "".concat(this.element.className, " dz-browser-not-supported");
        var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
        try {
            for(var _iterator = this.element.getElementsByTagName("div")[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                var child = _step.value;
                if (/(^| )dz-message($| )/.test(child.className)) {
                    messageElement = child;
                    child.className = "dz-message"; // Removes the 'dz-default' class
                    break;
                }
            }
        } catch (err) {
            _didIteratorError = true;
            _iteratorError = err;
        } finally{
            try {
                if (!_iteratorNormalCompletion && _iterator.return != null) {
                    _iterator.return();
                }
            } finally{
                if (_didIteratorError) {
                    throw _iteratorError;
                }
            }
        }
        if (!messageElement) {
            messageElement = (0, _dropzoneDefault.default).createElement('<div class="dz-message"><span></span></div>');
            this.element.appendChild(messageElement);
        }
        var span = messageElement.getElementsByTagName("span")[0];
        if (span) {
            if (span.textContent != null) span.textContent = this.options.dictFallbackMessage;
            else if (span.innerText != null) span.innerText = this.options.dictFallbackMessage;
        }
        return this.element.appendChild(this.getFallbackForm());
    },
    /**
   * Gets called to calculate the thumbnail dimensions.
   *
   * It gets `file`, `width` and `height` (both may be `null`) as parameters and must return an object containing:
   *
   *  - `srcWidth` & `srcHeight` (required)
   *  - `trgWidth` & `trgHeight` (required)
   *  - `srcX` & `srcY` (optional, default `0`)
   *  - `trgX` & `trgY` (optional, default `0`)
   *
   * Those values are going to be used by `ctx.drawImage()`.
   */ resize: function(file, width, height, resizeMethod) {
        var info = {
            srcX: 0,
            srcY: 0,
            srcWidth: file.width,
            srcHeight: file.height
        };
        var srcRatio = file.width / file.height;
        // Automatically calculate dimensions if not specified
        if (width == null && height == null) {
            width = info.srcWidth;
            height = info.srcHeight;
        } else if (width == null) width = height * srcRatio;
        else if (height == null) height = width / srcRatio;
        // Make sure images aren't upscaled
        width = Math.min(width, info.srcWidth);
        height = Math.min(height, info.srcHeight);
        var trgRatio = width / height;
        if (info.srcWidth > width || info.srcHeight > height) {
            // Image is bigger and needs rescaling
            if (resizeMethod === "crop") {
                if (srcRatio > trgRatio) {
                    info.srcHeight = file.height;
                    info.srcWidth = info.srcHeight * trgRatio;
                } else {
                    info.srcWidth = file.width;
                    info.srcHeight = info.srcWidth / trgRatio;
                }
            } else if (resizeMethod === "contain") {
                // Method 'contain'
                if (srcRatio > trgRatio) height = width / srcRatio;
                else width = height * srcRatio;
            } else throw new Error("Unknown resizeMethod '".concat(resizeMethod, "'"));
        }
        info.srcX = (file.width - info.srcWidth) / 2;
        info.srcY = (file.height - info.srcHeight) / 2;
        info.trgWidth = width;
        info.trgHeight = height;
        return info;
    },
    /**
   * Can be used to transform the file (for example, resize an image if necessary).
   *
   * The default implementation uses `resizeWidth` and `resizeHeight` (if provided) and resizes
   * images according to those dimensions.
   *
   * Gets the `file` as the first parameter, and a `done()` function as the second, that needs
   * to be invoked with the file when the transformation is done.
   */ transformFile: function(file, done) {
        if ((this.options.resizeWidth || this.options.resizeHeight) && file.type.match(/image.*/)) return this.resizeImage(file, this.options.resizeWidth, this.options.resizeHeight, this.options.resizeMethod, done);
        else return done(file);
    },
    /**
   * A string that contains the template used for each dropped
   * file. Change it to fulfill your needs but make sure to properly
   * provide all elements.
   *
   * If you want to use an actual HTML element instead of providing a String
   * as a config option, you could create a div with the id `tpl`,
   * put the template inside it and provide the element like this:
   *
   *     document
   *       .querySelector('#tpl')
   *       .innerHTML
   *
   */ previewTemplate: (0, _previewTemplateHtmlDefault.default),
    /*
   Those functions register themselves to the events on init and handle all
   the user interface specific stuff. Overwriting them won't break the upload
   but can break the way it's displayed.
   You can overwrite them if you don't like the default behavior. If you just
   want to add an additional event handler, register it on the dropzone object
   and don't overwrite those options.
   */ // Those are self explanatory and simply concern the DragnDrop.
    drop: function(e) {
        return this.element.classList.remove("dz-drag-hover");
    },
    dragstart: function(e) {},
    dragend: function(e) {
        return this.element.classList.remove("dz-drag-hover");
    },
    dragenter: function(e) {
        return this.element.classList.add("dz-drag-hover");
    },
    dragover: function(e) {
        return this.element.classList.add("dz-drag-hover");
    },
    dragleave: function(e) {
        return this.element.classList.remove("dz-drag-hover");
    },
    paste: function(e) {},
    // Called whenever there are no files left in the dropzone anymore, and the
    // dropzone should be displayed as if in the initial state.
    reset: function() {
        return this.element.classList.remove("dz-started");
    },
    // Called when a file is added to the queue
    // Receives `file`
    addedfile: function(file) {
        var _this = this;
        if (this.element === this.previewsContainer) this.element.classList.add("dz-started");
        if (this.previewsContainer && !this.options.disablePreviews) {
            file.previewElement = (0, _dropzoneDefault.default).createElement(this.options.previewTemplate.trim());
            file.previewTemplate = file.previewElement; // Backwards compatibility
            this.previewsContainer.appendChild(file.previewElement);
            var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
            try {
                for(var _iterator = file.previewElement.querySelectorAll("[data-dz-name]")[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                    var node = _step.value;
                    node.textContent = file.name;
                }
            } catch (err) {
                _didIteratorError = true;
                _iteratorError = err;
            } finally{
                try {
                    if (!_iteratorNormalCompletion && _iterator.return != null) {
                        _iterator.return();
                    }
                } finally{
                    if (_didIteratorError) {
                        throw _iteratorError;
                    }
                }
            }
            var _iteratorNormalCompletion1 = true, _didIteratorError1 = false, _iteratorError1 = undefined;
            try {
                for(var _iterator1 = file.previewElement.querySelectorAll("[data-dz-size]")[Symbol.iterator](), _step1; !(_iteratorNormalCompletion1 = (_step1 = _iterator1.next()).done); _iteratorNormalCompletion1 = true){
                    node = _step1.value;
                    node.innerHTML = this.filesize(file.size);
                }
            } catch (err) {
                _didIteratorError1 = true;
                _iteratorError1 = err;
            } finally{
                try {
                    if (!_iteratorNormalCompletion1 && _iterator1.return != null) {
                        _iterator1.return();
                    }
                } finally{
                    if (_didIteratorError1) {
                        throw _iteratorError1;
                    }
                }
            }
            if (this.options.addRemoveLinks) {
                file._removeLink = (0, _dropzoneDefault.default).createElement('<a class="dz-remove" href="javascript:undefined;" data-dz-remove>'.concat(this.options.dictRemoveFile, "</a>"));
                file.previewElement.appendChild(file._removeLink);
            }
            var removeFileEvent = function(e) {
                e.preventDefault();
                e.stopPropagation();
                if (file.status === (0, _dropzoneDefault.default).UPLOADING) return (0, _dropzoneDefault.default).confirm(_this.options.dictCancelUploadConfirmation, function() {
                    return _this.removeFile(file);
                });
                else {
                    if (_this.options.dictRemoveFileConfirmation) return (0, _dropzoneDefault.default).confirm(_this.options.dictRemoveFileConfirmation, function() {
                        return _this.removeFile(file);
                    });
                    else return _this.removeFile(file);
                }
            };
            var _iteratorNormalCompletion2 = true, _didIteratorError2 = false, _iteratorError2 = undefined;
            try {
                for(var _iterator2 = file.previewElement.querySelectorAll("[data-dz-remove]")[Symbol.iterator](), _step2; !(_iteratorNormalCompletion2 = (_step2 = _iterator2.next()).done); _iteratorNormalCompletion2 = true){
                    var removeLink = _step2.value;
                    removeLink.addEventListener("click", removeFileEvent);
                }
            } catch (err) {
                _didIteratorError2 = true;
                _iteratorError2 = err;
            } finally{
                try {
                    if (!_iteratorNormalCompletion2 && _iterator2.return != null) {
                        _iterator2.return();
                    }
                } finally{
                    if (_didIteratorError2) {
                        throw _iteratorError2;
                    }
                }
            }
        }
    },
    // Called whenever a file is removed.
    removedfile: function(file) {
        if (file.previewElement != null && file.previewElement.parentNode != null) file.previewElement.parentNode.removeChild(file.previewElement);
        return this._updateMaxFilesReachedClass();
    },
    // Called when a thumbnail has been generated
    // Receives `file` and `dataUrl`
    thumbnail: function(file, dataUrl) {
        if (file.previewElement) {
            file.previewElement.classList.remove("dz-file-preview");
            var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
            try {
                for(var _iterator = file.previewElement.querySelectorAll("[data-dz-thumbnail]")[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                    var thumbnailElement = _step.value;
                    thumbnailElement.alt = file.name;
                    thumbnailElement.src = dataUrl;
                }
            } catch (err) {
                _didIteratorError = true;
                _iteratorError = err;
            } finally{
                try {
                    if (!_iteratorNormalCompletion && _iterator.return != null) {
                        _iterator.return();
                    }
                } finally{
                    if (_didIteratorError) {
                        throw _iteratorError;
                    }
                }
            }
            return setTimeout(function() {
                return file.previewElement.classList.add("dz-image-preview");
            }, 1);
        }
    },
    // Called whenever an error occurs
    // Receives `file` and `message`
    error: function(file, message) {
        if (file.previewElement) {
            file.previewElement.classList.add("dz-error");
            if (typeof message !== "string" && message.error) message = message.error;
            var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
            try {
                for(var _iterator = file.previewElement.querySelectorAll("[data-dz-errormessage]")[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                    var node = _step.value;
                    node.textContent = message;
                }
            } catch (err) {
                _didIteratorError = true;
                _iteratorError = err;
            } finally{
                try {
                    if (!_iteratorNormalCompletion && _iterator.return != null) {
                        _iterator.return();
                    }
                } finally{
                    if (_didIteratorError) {
                        throw _iteratorError;
                    }
                }
            }
        }
    },
    errormultiple: function() {},
    // Called when a file gets processed. Since there is a queue, not all added
    // files are processed immediately.
    // Receives `file`
    processing: function(file) {
        if (file.previewElement) {
            file.previewElement.classList.add("dz-processing");
            if (file._removeLink) return file._removeLink.innerHTML = this.options.dictCancelUpload;
        }
    },
    processingmultiple: function() {},
    // Called whenever the upload progress gets updated.
    // Receives `file`, `progress` (percentage 0-100) and `bytesSent`.
    // To get the total number of bytes of the file, use `file.size`
    uploadprogress: function(file, progress, bytesSent) {
        var _iteratorNormalCompletion = true, _didIteratorError = false, _iteratorError = undefined;
        if (file.previewElement) try {
            for(var _iterator = file.previewElement.querySelectorAll("[data-dz-uploadprogress]")[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true){
                var node = _step.value;
                node.nodeName === "PROGRESS" ? node.value = progress : node.style.width = "".concat(progress, "%");
            }
        } catch (err) {
            _didIteratorError = true;
            _iteratorError = err;
        } finally{
            try {
                if (!_iteratorNormalCompletion && _iterator.return != null) {
                    _iterator.return();
                }
            } finally{
                if (_didIteratorError) {
                    throw _iteratorError;
                }
            }
        }
    },
    // Called whenever the total upload progress gets updated.
    // Called with totalUploadProgress (0-100), totalBytes and totalBytesSent
    totaluploadprogress: function() {},
    // Called just before the file is sent. Gets the `xhr` object as second
    // parameter, so you can modify it (for example to add a CSRF token) and a
    // `formData` object to add additional information.
    sending: function() {},
    sendingmultiple: function() {},
    // When the complete upload is finished and successful
    // Receives `file`
    success: function(file) {
        if (file.previewElement) return file.previewElement.classList.add("dz-success");
    },
    successmultiple: function() {},
    // When the upload is canceled.
    canceled: function(file) {
        return this.emit("error", file, this.options.dictUploadCanceled);
    },
    canceledmultiple: function() {},
    // When the upload is finished, either with success or an error.
    // Receives `file`
    complete: function(file) {
        if (file._removeLink) file._removeLink.innerHTML = this.options.dictRemoveFile;
        if (file.previewElement) return file.previewElement.classList.add("dz-complete");
    },
    completemultiple: function() {},
    maxfilesexceeded: function() {},
    maxfilesreached: function() {},
    queuecomplete: function() {},
    addedfiles: function() {}
};
exports.default = defaultOptions;

},{"./dropzone":"1u0OZ","bundle-text:./preview-template.html":"gUxwX","@parcel/transformer-js/src/esmodule-helpers.js":"3Qwoy"}],"gUxwX":[function(require,module,exports) {
module.exports = "<div class=\"dz-preview dz-file-preview\">\n  <div class=\"dz-image\"><img data-dz-thumbnail=\"\"></div>\n  <div class=\"dz-details\">\n    <div class=\"dz-size\"><span data-dz-size=\"\"></span></div>\n    <div class=\"dz-filename\"><span data-dz-name=\"\"></span></div>\n  </div>\n  <div class=\"dz-progress\">\n    <span class=\"dz-upload\" data-dz-uploadprogress=\"\"></span>\n  </div>\n  <div class=\"dz-error-message\"><span data-dz-errormessage=\"\"></span></div>\n  <div class=\"dz-success-mark\">\n    <svg width=\"54\" height=\"54\" viewbox=\"0 0 54 54\" fill=\"white\" xmlns=\"http://www.w3.org/2000/svg\">\n      <path d=\"M10.2071 29.7929L14.2929 25.7071C14.6834 25.3166 15.3166 25.3166 15.7071 25.7071L21.2929 31.2929C21.6834 31.6834 22.3166 31.6834 22.7071 31.2929L38.2929 15.7071C38.6834 15.3166 39.3166 15.3166 39.7071 15.7071L43.7929 19.7929C44.1834 20.1834 44.1834 20.8166 43.7929 21.2071L22.7071 42.2929C22.3166 42.6834 21.6834 42.6834 21.2929 42.2929L10.2071 31.2071C9.81658 30.8166 9.81658 30.1834 10.2071 29.7929Z\"></path>\n    </svg>\n  </div>\n  <div class=\"dz-error-mark\">\n    <svg width=\"54\" height=\"54\" viewbox=\"0 0 54 54\" fill=\"white\" xmlns=\"http://www.w3.org/2000/svg\">\n      <path d=\"M26.2929 20.2929L19.2071 13.2071C18.8166 12.8166 18.1834 12.8166 17.7929 13.2071L13.2071 17.7929C12.8166 18.1834 12.8166 18.8166 13.2071 19.2071L20.2929 26.2929C20.6834 26.6834 20.6834 27.3166 20.2929 27.7071L13.2071 34.7929C12.8166 35.1834 12.8166 35.8166 13.2071 36.2071L17.7929 40.7929C18.1834 41.1834 18.8166 41.1834 19.2071 40.7929L26.2929 33.7071C26.6834 33.3166 27.3166 33.3166 27.7071 33.7071L34.7929 40.7929C35.1834 41.1834 35.8166 41.1834 36.2071 40.7929L40.7929 36.2071C41.1834 35.8166 41.1834 35.1834 40.7929 34.7929L33.7071 27.7071C33.3166 27.3166 33.3166 26.6834 33.7071 26.2929L40.7929 19.2071C41.1834 18.8166 41.1834 18.1834 40.7929 17.7929L36.2071 13.2071C35.8166 12.8166 35.1834 12.8166 34.7929 13.2071L27.7071 20.2929C27.3166 20.6834 26.6834 20.6834 26.2929 20.2929Z\"></path>\n    </svg>\n  </div>\n</div>\n<script src=\"/preview-template.3b3362d3.js\"></script>";

},{}]},["hq6rc","42jTR"], "42jTR", "parcelRequire1768")

//# sourceMappingURL=dropzone-min.js.map
