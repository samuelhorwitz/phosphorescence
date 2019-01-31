# Web application

This is the code for the actual web application. It's all written in Vue JS.

Note that we have a separate origin to actually run set-building. This is because we want to isolate user-generated code from the main origin which houses credentials and other things we don't want people to XSS attack. This separate origin is Eos and you will find it in the Eos folder along with all the actual playlist-builder running code.