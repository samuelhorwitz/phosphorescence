# Web application

This is the code for the actual web application. The frontend is written in Vue JS and the API is written in Go.

Note that we have a separate origin to actually run set-building. This is because we want to isolate user-generated code from the main origin which houses credentials and other things we don't want people to XSS attack. This separate origin is Eos and you will find it in the Eos folder along with all the actual playlist-builder running code.

Phosphorescence's "vivid sunrise" background image and branding is a derivative work of a photograph by Karl Magnuson which can be found here: https://unsplash.com/photos/HQR_JXd-fPs

Many thanks to Karl Magnuson for this freely available photograph.