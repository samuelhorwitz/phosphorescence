import Logger from './logger';
import Queue from './queue';

const messageLoop = Symbol('messageLoop');
const channelKnockBuilder = Symbol('channelKnockBuilder');
const postAndWait = Symbol('postAndWait');
const postMessage = Symbol('postMessage');
const postMessageAfterInitialized = Symbol('postMessageAfterInitialized');
const wrapMessage = Symbol('wrapMessage');
const isWrapped = Symbol('isWrapped');
const newDestroyablePromise = Symbol('newDestroyablePromise');
const createOneTimePort = Symbol('createOneTimePort');
const createMultiusePort = Symbol('createMultiusePort');
const createOneTimeChannel = Symbol('createOneTimeChannel');
const checkKnocker = Symbol('checkKnocker');
const checkListener = Symbol('checkListener');
const listenIframe = Symbol('listenIframe');
const listenWorker = Symbol('listenWorker');
const LISTENER_TYPE = Symbol('listener');
const KNOCKER_TYPE = Symbol('knocker');

export default class SecureMessenger {
    #logger = new Logger(location.origin, process.env.NODE_ENV === 'production' ? 'production' : 'development');
    #initializerCalled = false;
    #remoteOrigin;
    #closed = false;
    #queue;
    #destructionQueue = [];
    #type;
    #initializationComplete;

    constructor(remoteOrigin) {
        if (remoteOrigin === '*') {
            this.#logger.warn('Disallowing insecure origin "*"');
            throw new Error('Cannot use insecure origin "*"');
        }
        this.#logger.debug(`Secure messenger initialized for origin ${remoteOrigin}`);
        this.#queue = new Queue(this.#logger);
        this.#remoteOrigin = remoteOrigin;
    }

    listen(toListenOn) {
        if (this.#initializerCalled) {
            throw new Error('Already initialized');
        }
        if (!toListenOn || (toListenOn instanceof HTMLIFrameElement)) {
            return this[listenIframe](toListenOn);
        } else if (toListenOn instanceof Worker) {
            return this[listenWorker](toListenOn);
        }
        throw new Error('Invalid listener element');
    }

    [listenIframe](iframe) {
        if (iframe && !(iframe instanceof HTMLIFrameElement)) {
            throw new Error('Invalid iframe element');
        }
        this.#type = LISTENER_TYPE;
        this.#logger.debug('Initializing listener (iframe)');
        this.#initializerCalled = true;
        return this.#initializationComplete = this[newDestroyablePromise]((resolve, reject) => {
            let cleanedUp = false;
            let messageHandler = async event => {
                if (event.origin !== this.#remoteOrigin) {
                    return;
                }
                this.#logger.debug('Listener received response', event);
                if (event.ports.length == 0 || !this[checkKnocker](event.data)) {
                    if (event.data.$version !== 1) {
                        this.#logger.warn('Listener initialization failed due to invalid version', event);
                    }
                    else {
                        this.#logger.warn('Listener initialization failed due to invalid response', event);
                    }
                    reject(event);
                    finishListeningOnWindow(event);
                    return;
                }
                if (iframe && event.source !== iframe.contentWindow) {
                    this.#logger.debug('Incorrect iframe source, will wait for next event', iframe, event);
                    return;
                }
                finishListeningOnWindow(event);
                try {
                    this.#logger.debug('Listener sending handshake response');
                    let handshakePromise = this[postMessage]({$type: '🤝', $version: 1});
                    this.#logger.debug('Listener initializing message loop');
                    this[messageLoop](event.data.$responsePort);
                    this.#logger.debug('Listener waiting for handshake third-part response');
                    let {data} = await handshakePromise;
                    this.#logger.debug('Listener received handshake third-part', data);
                    if (data.$type !== '👂') {
                        throw new Error('Knocker failed to properly complete three-way handshake');
                    }
                } catch (e) {
                    this.#logger.warn('Listener initialization failed in response handling', e);
                    reject(e);
                    return;
                }
                resolve();
            };
            let finishListeningOnWindow = event => {
                event.stopImmediatePropagation();
                cleanupWindowListeners();
            };
            let messageErrorHandler = event => {
                if (event.origin !== this.#remoteOrigin) {
                    return;
                }
                if (iframe && event.source !== iframe.contentWindow) {
                    return;
                }
                this.#logger.warn('Listener initialization failed due to errored message', event);
                finishListeningOnWindow(event);
                reject(event);
            };
            let cleanupWindowListeners = () => {
                if (cleanedUp) {
                    this.#logger.debug('Temporary window listeners already cleaned up');
                    return;
                }
                this.#logger.debug('Listener cleaning up temporary window listeners');
                removeEventListener('message', messageHandler);
                removeEventListener('messageerror', messageErrorHandler);
                cleanedUp = true;
            };
            this.#logger.debug('Listener adding temporary window listeners');
            addEventListener('message', messageHandler);
            addEventListener('messageerror', messageErrorHandler);
            this.#destructionQueue.push(cleanupWindowListeners);
        });
    }

    [listenWorker](worker) {
        if (!(worker instanceof Worker)) {
            throw new Error('Invalid worker');
        }
        if (this.#remoteOrigin !== location.origin) {
            throw new Error('Origin must be current origin if worker listener');
        }
        this.#type = LISTENER_TYPE;
        this.#logger.debug('Initializing listener (worker)');
        this.#initializerCalled = true;
        return this.#initializationComplete = this[newDestroyablePromise]((resolve, reject) => {
            let cleanedUp = false;
            let messageHandler = async event => {
                this.#logger.debug('Listener received response', event);
                if (event.ports.length == 0 || !this[checkKnocker](event.data)) {
                    if (event.data.$version !== 1) {
                        this.#logger.warn('Listener initialization failed due to invalid version', event);
                    }
                    else {
                        this.#logger.warn('Listener initialization failed due to invalid response', event);
                    }
                    reject(event);
                    finishListeningOnWindow(event);
                    return;
                }
                finishListeningOnWindow(event);
                try {
                    this.#logger.debug('Listener sending handshake response');
                    let handshakePromise = this[postMessage]({$type: '🤝', $version: 1});
                    this.#logger.debug('Listener initializing message loop');
                    this[messageLoop](event.data.$responsePort);
                    this.#logger.debug('Listener waiting for handshake third-part response');
                    let {data} = await handshakePromise;
                    this.#logger.debug('Listener received handshake third-part', data);
                    if (data.$type !== '👂') {
                        throw new Error('Knocker failed to properly complete three-way handshake');
                    }
                } catch (e) {
                    this.#logger.warn('Listener initialization failed in response handling', e);
                    reject(e);
                    return;
                }
                resolve();
            };
            let finishListeningOnWindow = event => {
                event.stopImmediatePropagation();
                cleanupWindowListeners();
            };
            let messageErrorHandler = event => {
                this.#logger.warn('Listener initialization failed due to errored message', event);
                finishListeningOnWindow(event);
                reject(event);
            };
            let cleanupWindowListeners = () => {
                if (cleanedUp) {
                    this.#logger.debug('Temporary worker listeners already cleaned up');
                    return;
                }
                this.#logger.debug('Listener cleaning up temporary worker listeners');
                worker.removeEventListener('message', messageHandler);
                worker.removeEventListener('messageerror', messageErrorHandler);
                cleanedUp = true;
            };
            this.#logger.debug('Listener adding temporary worker listeners');
            worker.addEventListener('message', messageHandler);
            worker.addEventListener('messageerror', messageErrorHandler);
            this.#destructionQueue.push(cleanupWindowListeners);
        });
    }

    knock(elements) {
        if (this.#initializerCalled) {
            throw new Error('Already initialized');
        }
        this.#type = KNOCKER_TYPE;
        this.#logger.debug('Initializing knocker');
        this.#initializerCalled = true;
        return this.#initializationComplete = this[newDestroyablePromise]((resolve, reject) => {
            let getPort = this[channelKnockBuilder](this[checkListener], event => {
                this.#logger.debug('Knock response received', event);
                this.#logger.debug('Knocker initializing message loop');
                this[messageLoop](event.data.$responsePort);
                resolve();
            }, event => {
                this.#logger.warn('Knock response error', event);
                reject(event);
            });
            this.#logger.debug(`Preparing to knock. You may see errors such as "Failed to execute 'postMessage' on 'DOMWindow'". These may be disregarded.`);
            for (let i = 0; i < elements.length; i++) {
                this.#logger.debug(`Knocking on element ${i}...`);
                let port = getPort();
                if ((typeof Worker !== 'undefined' && elements[i] instanceof Worker) || (typeof WorkerGlobalScope !== 'undefined' && elements[i] instanceof WorkerGlobalScope)) {
                    if (this.#remoteOrigin !== location.origin) {
                        throw new Error('Origin must be current origin if worker knock');
                    }
                    elements[i].postMessage({$type: '👋', $version: 1, $responsePort: port}, [port]);
                } else {
                    elements[i].postMessage({$type: '👋', $version: 1, $responsePort: port}, this.#remoteOrigin, [port]);
                }
            }
        });
    }

    messageHandlerLoop(handlerFn, errorHandlerFn) {
        if (this.#type === LISTENER_TYPE) {
            throw new Error('Only the knocker may initialize a message handler loop');
        }
        let keepLooping = true;
        let finalMessage = false;
        let finisher = () => {
            this.#logger.debug('Message handler loop will break after this message');
            finalMessage = true;
        };
        (async () => {
            this.#logger.debug('Message handler loop will wait for first two phases of handshake');
            await this.#initializationComplete;
            this.#logger.debug('Message handler loop ready to proceed after first two phases of handshake');
            let msg = {$type: '👂', $raw: true};
            while (keepLooping) {
                this.#logger.debug('Message handler loop next message');
                let eventData;
                try {
                    this.#logger.debug('Message handler loop posting message', msg);
                    eventData = await this[postMessageAfterInitialized](msg);
                    this.#logger.debug('Message handler loop response received', eventData);
                }
                catch (e) {
                    try {
                        this.#logger.warn('Message handler loop post and wait failed', e);
                        msg = await errorHandlerFn(e);
                        if (msg) {
                            this.#logger.debug('Message handler loop error handler responded with next message, continue loop', msg);
                            continue;
                        } else {
                            this.#logger.warn('Message handler loop error handler did not respond with message, break loop');
                            break;
                        }
                    } catch (e) {
                        this.#logger.warn('Message handler loop could not handle error, breaking loop', e);
                        break;
                    }
                }
                if (finalMessage) {
                    break;
                }
                this.#logger.debug('Message handler loop handling response', eventData);
                msg = this[wrapMessage](await handlerFn(eventData.data, eventData.interruptPort, finisher));
                this.#logger.debug('Message handler loop response handled, next message', msg);
            }
            this.#logger.debug('Message handler loop exited');
        })();
        return () => {
            this.#logger.debug('Message handler loop will not loop again');
            keepLooping = false;
        };
    }

    async finishHandshake() {
        if (this.#type === LISTENER_TYPE) {
            throw new Error('The knocker is responsible for finishing the three-way handshake');
        }
        this.#logger.debug('Finish handshake will wait for first two phases of handshake');
        await this.#initializationComplete;
        this.#logger.debug('Finish handshake ready to proceed after first two phases of handshake');
        return this[postMessageAfterInitialized]({$type: '👂', $raw: true});
    }

    postMessage(message) {
        return this[postMessageAfterInitialized](this[wrapMessage](message));
    }

    async openOneTimeInterruptListenerPort(message) {
        await this.#initializationComplete;
        this.#logger.debug('Opening one time interrupt port', message);
        let interruptResolver, interruptRejecter;
        let interruption = this[newDestroyablePromise]((resolve, reject) => {
            interruptResolver = resolve;
            interruptRejecter = reject;
        });
        return this[newDestroyablePromise]((resolve, reject) => {
            let port = this[createOneTimePort](event => {
                this.#logger.debug('One time interrupt triggered', event);
                interruptResolver({data: event.data});
            }, event => {
                this.#logger.warn('One time interrupt failed', event);
                interruptRejecter(event);
            });
            this.#queue.push({
                message: this[wrapMessage](message, port),
                handler: (data) => {
                    this.#logger.debug('Message response received for one time interrupt port', data);
                    resolve({data, interruption});
                },
                errorHandler: error => {
                    this.#logger.warn('Message response for one time interrupt port errored', error);
                    reject(error);
                }
            });
        });
    }

    async openInterruptListenerPort(message, handlerFn, errorHandlerFn) {
        await this.#initializationComplete;
        this.#logger.debug('Opening interrupt port', message);
        return this[newDestroyablePromise]((resolve, reject) => {
            let {port, closer} = this[createMultiusePort](event => {
                this.#logger.debug('Interrupt triggered', event);
                handlerFn(event.data);
            }, event => {
                this.#logger.warn('Interrupt failed', event);
                errorHandlerFn(event);
            });
            this.#queue.push({
                message: this[wrapMessage](message, port),
                handler: (data) => {
                    this.#logger.debug('Message response received for interrupt port', data);
                    resolve({data, closer});
                },
                errorHandler: error => {
                    this.#logger.warn('Message response for interrupt port errored', error);
                    reject(error);
                }
            });
        });
    }

    async respondWithInterruptListenerPort(message, handlerFn, errorHandlerFn) {
        await this.#initializationComplete;
        this.#logger.debug('Opening interrupt port', message);
        let {port, closer} = this[createMultiusePort](event => {
            this.#logger.debug('Interrupt triggered', event);
            handlerFn(event.data);
        }, event => {
            this.#logger.warn('Interrupt failed', event);
            errorHandlerFn(event);
        });
        return this[wrapMessage](message, port);
    }

    close() {
        if (this.#closed) {
            throw new Error('Channel already closed');
        }
        this.#logger.debug('Closing channel');
        this.#closed = true;
        this.#queue.destroy();
        this.#destructionQueue.forEach((destroyer, i) => {
            this.#logger.debug(`Destroying destroyable ${i}...`);
            destroyer();
        });
        this.#destructionQueue.length = 0;
    }

    [wrapMessage](message, interruptPort) {
        if (message && message[isWrapped]) {
            this.#logger.debug('Message already wrapped', message);
            return message;
        }
        let wrappedMessage = {[isWrapped]: true, message};
        if (interruptPort) {
            wrappedMessage.$interruptPort = interruptPort;
        }
        return wrappedMessage;
    }

    async [postMessageAfterInitialized](message) {
        await this.#initializationComplete;
        return this[postMessage](message);
    }

    [postMessage](message) {
        this.#logger.debug('Posting message', message);
        return this[newDestroyablePromise]((resolve, reject) => {
            this.#queue.push({
                message,
                handler: (data, interruptPort) => {
                    this.#logger.debug('Message response received', data, interruptPort);
                    resolve({data, interruptPort});
                },
                errorHandler: error => {
                    this.#logger.warn('Message response errored', error);
                    reject(error);
                }
            });
        });
    }

    [channelKnockBuilder](checkFn, handleFn, handleErrorFn) {
        this.#logger.debug('Preparing knocker');
        let channels = [];
        let destroyChannels = () => {
            if (channels.length == 0) {
                this.#logger.debug('No knock channels left to destroy');
                return;
            }
            this.#logger.debug('Destroying knock channels');
            channels.forEach((chan, i) => {
                this.#logger.debug(`Destroying knock channel ${i}...`, chan.port1);
                chan.port1.close();
            });
            channels.length = 0;
        };
        this.#destructionQueue.push(destroyChannels);
        return () => {
            this.#logger.debug('Building a knock channel');
            let chan = this[createOneTimeChannel](event => {
                this.#logger.debug('Knock response found', event);
                destroyChannels();
                if (event.ports.length > 0 && checkFn(event.data)) {
                    this.#logger.debug('Knock response passes checks', event);
                    handleFn(event);
                }
                else {
                    this.#logger.warn('Knock response is invalid, erroring', event);
                    handleErrorFn(event);
                }
            }, event => {
                this.#logger.warn('Knock attempt failed', event);
                handleErrorFn(event);
            });
            channels.push(chan);
            return chan.port2;
        };
    }

    async [messageLoop](initialPort) {
        let port = initialPort;
        this.#logger.debug('Message loop starting');
        while (!this.#closed) {
            this.#logger.debug('Next message loop');
            let queueValue;
            try {
                this.#logger.debug('Message loop wants to shift message off queue');
                queueValue = await this.#queue.shift();
            }
            catch (e) {
                if (this.#closed) {
                    this.#logger.debug('Message loop queue shift failed but channel closed, so failure expected', e);
                    break;
                }
                this.#logger.warn('Message loop shift from queue failed', e);
                throw new Error(e);
            }
            this.#logger.debug('Message loop got message from queue', queueValue);
            let {message, handler, errorHandler} = queueValue;
            let event;
            try {
                this.#logger.debug('Message loop wants to post message and wait for response', queueValue);
                event = await this[postAndWait](port, message);
                this.#logger.debug('Message loop received response', event);
                port = event.data.$responsePort;
            }
            catch (e) {
                this.#logger.warn('Message loop request failed and cannot recover', e);
                await errorHandler(e);
                this.#logger.debug('Message loop error handler completed', e);
                break;
            }
            try {
                this.#logger.debug('Message loop forwarding event to handler', event);
                let handlerMessage;
                if (event.data.$raw) {
                    handlerMessage = event.data;
                }
                else {
                    handlerMessage = event.data.message;
                }
                await handler(handlerMessage, event.data.$interruptPort);
                this.#logger.debug('Message loop handler completed', event);
            }
            catch (e) {
                this.#logger.warn('Message loop handler failed', e);
                await errorHandler(e);
                this.#logger.debug('Message loop error handler completed', e);
            }
        }
        this.#logger.debug('Message loop exited');
    }

    [postAndWait](port, message) {
        return this[newDestroyablePromise]((resolve, reject) => {
            let transferList = [];
            if (message.$interruptPort) {
                transferList.unshift(message.$interruptPort);
            }
            let newMessage = Object.assign({}, message, {$responsePort: this[createOneTimePort](resolve, reject)})
            transferList.unshift(newMessage.$responsePort);
            port.postMessage(newMessage, transferList);
        });
    }

    [newDestroyablePromise](fn) {
        return new Promise((resolve, reject) => {
            if (this.#closed) {
                reject('Channel closed');
                return;
            }
            this.#destructionQueue.push(() => reject('Channel closed'));
            fn(resolve, reject);
        });
    }

    [createOneTimeChannel](handler, errorHandler) {
        if (this.#closed) {
            errorHandler('Channel closed');
            return;
        }
        let {chan, closer} = _createOneTimeChannel(handler, errorHandler);
        this.#destructionQueue.push(closer);
        return chan;
    }

    [createOneTimePort](handler, errorHandler) {
        if (this.#closed) {
            errorHandler('Channel closed');
            return;
        }
        let {port, closer} = _createOneTimePort(handler, errorHandler);
        this.#destructionQueue.push(closer);
        return port;
    }

    [createMultiusePort](handler, errorHandler) {
        if (this.#closed) {
            errorHandler('Channel closed');
            return;
        }
        let {port, closer} = _createMultiusePort(handler, errorHandler);
        this.#destructionQueue.push(closer);
        return {port, closer};
    }

    [checkListener]({$type, $version}) {
        if ($version !== 1) {
            this.#logger.warn('Invalid listener version', $version);
        }
        return $type === '🤝' && $version === 1;
    }

    [checkKnocker]({$type, $version}) {
        if ($version !== 1) {
            this.#logger.warn('Invalid knocker version', $version);
        }
        return $type === '👋' && $version === 1;
    }
}

function _createOneTimeChannel(handler, errorHandler) {
    let chan = new MessageChannel();
    let handleAndClose = event => {
        chan.port1.close();
        handler(event);
    };
    let handleErrorAndClose = event => {
        chan.port1.close();
        errorHandler(event);
    };
    chan.port1.addEventListener('message', handleAndClose, {once: true});
    chan.port1.addEventListener('messageerror', handleErrorAndClose, {once: true});
    chan.port1.start();
    return {chan, closer: () => {
        chan.port1.close();
        chan.port1.removeEventListener('message', handleAndClose);
        chan.port1.removeEventListener('messageerror', handleErrorAndClose);
    }};
}

function _createOneTimePort(handler, errorHandler) {
    let {chan, closer} = _createOneTimeChannel(handler, errorHandler);
    return {port: chan.port2, closer};
}

function _createMultiuseChannel(handler, errorHandler) {
    let chan = new MessageChannel();
    chan.port1.addEventListener('message', handler);
    chan.port1.addEventListener('messageerror', errorHandler);
    chan.port1.start();
    return {chan, closer: () => {
        chan.port1.close();
        chan.port1.removeEventListener('message', handler);
        chan.port1.removeEventListener('messageerror', errorHandler);
    }};
}

function _createMultiusePort(handler, errorHandler) {
    let {chan, closer} = _createMultiuseChannel(handler, errorHandler);
    return {port: chan.port2, closer};
}