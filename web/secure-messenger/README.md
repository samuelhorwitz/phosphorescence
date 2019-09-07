# Secure Messenger

The secure messenger is a wrapper around browser `postMessage` communication that attempts to serialize communication through the use of one time ports and an enforced listen and respond protocol. It also supports interrupt messages for secondary notification-based one-way communication.

Here is how I believe it would be explained in the OSI model:

- Layers 1-2 (Physical, Datalink): Not really relevant but I guess technically provided by the browser.
- Layer 3 (Network): This is provided by the browser's builtin `postMessage` cross-origin event system.
- Layer 4 (Transport): This is the layer the secure messenger provides.
- Layer 5 (Session): Secure messenger also touches some of this layer, regarding connection session not application session
- Layers 6-7 (Presentation, Application): This is Phosphorescence and Eos and the application-level information being transported around between them.

Communicating through raw `postMessage` and `onmessage` handlers is doable and obviously simpler for basic applications. However, secure messenger intends to take a more controlled approach. Certain forms of `postMessage` in the browser have built-in security features such as enforced origin (for iframe messaging). We want to build on the existing features.

Browsers now also support `postMessage` not just between window or context objects (iframes, popups, workers, etc.) but on a new concept called `MessageChannel`s. These message channels are created and provide two ports (which are "transferable objects") which can be listened to and emitted on. Transferable objects are objects who can be sent along in a `postMessage` and their owner moves from one window or context to another.

The easiest way to think about a message channel is that it sets up a children's soup can and string "telephone". One soup can is passed to another context and now both sides can talk into their side's can and put their ear towards it to listen.

The secure messenger is based around the notion that communication should occur in a tightly controlled loop where every message includes a new one-time-use soup can. Messages may not be sent out of order, A will message and wait for B to receive and respond before A can message again. This waiting is all abstracted away from the application using async queues. The technical term I guess is "half-duplex".

Sometimes, there are things outside the scope of typical protocol communication and cannot be architected in a send and wait loop. These are considered "interrupts" and are usually user interaction driven.

It is possible using the secure messenger to set up extra channels to handle interrupts. Interrupts are one-way only ("simplex"). There is the ability to set up multi-use interrupts and one-time interrupts. The host who wishes to listen for an interrupt sends down a request. The remote must have application-level code to handle this request but the creation of ports is all handled by secure messenger. The remote may now hold onto this port and use it to emit interrupts as needed and the host can now hear them. The host will not send anything down these ports as they are listen only. The wrapping of the interrupt port is prettier than working with raw ports directly: one time ports get a promise back that can be `await`ed on and multiuse ports are as simple as passing in a handler and error handler function.

Ports can be used by sending messages one-off and waiting for response (promises are used) or intializing a message handler loop (which can remove some of the boilerplate fatigue for constant communication application protocols).

To initialize a port, the secure messenger is set up on both contexts. There must be some level of out of band application protocol agreement, as with all `postMessage` communication, but basically one origin "knocks" and the other "listens". This uses window context messaging. Once the knocker and listen are in agreement, the communication may begin. The knocker will initialize a handler loop using the easy message loop function or by manually sending the third part of the three-way handshake acknowledging that it's listening and manually handling responses.

## Security model and limitations
The secure messenger's security is subtle and usage does not imply magical security all around. The role of the secure messenger is purely to provide a maintained connection between two origins through a gateway on each origin where all messages are strictly ordered and every message is sent through a one time use port and each origin must wait for a response between every sent message. It also provides a simplex interrupt alternate channel purely for event bubbling.

This should allow for origins to not have to have open message listeners scattered in the code and takes care of accounting for low level communication ordering. Application level ordering and state is still the responsibility of the application. Untrusted code running on the same origin as the gateway could easily execute a MITM attack, as is always the case with untrusted code being run. There is no encryption, the browser's cross-origin messaging is assumed to be secure because we are assuming the user's device is secure but obviously if it's not there are much bigger issues.

Keeping the one time ports "secret" by the nature of them being referenced only inside a "safe" closure helps prevent other code from being able to post messages maliciously. If both sides only communicate through these secure channels, then some form of XSS or whatever that manages to run shouldn't be able to just use the global `postMessage` directly as no one will be listening to that (as they will only listen on the currently active port that should be closured away from global use).