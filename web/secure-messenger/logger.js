function s4() {
    return Math.floor((1 + Math.random()) * 0x10000).toString(16).substring(1);
}

function uid() {
    return s4() + s4();
}

var colorsUnshuffled = ['MediumVioletRed', 'Crimson', 'OrangeRed', 'DarkKhaki', 'RosyBrown', 'Brown', 'SeaGreen', 'DarkCyan', 'MediumBlue', 'MediumOrchid', 'MediumPurple', 'DarkSlateBlue'];

function shuffle(array) {
    var currentIndex = array.length, temporaryValue, randomIndex;
    // While there remain elements to shuffle...
    while (0 !== currentIndex) {
        // Pick a remaining element...
        randomIndex = Math.floor(Math.random() * currentIndex);
        currentIndex -= 1;
        // And swap it with the current element.
        temporaryValue = array[currentIndex];
        array[currentIndex] = array[randomIndex];
        array[randomIndex] = temporaryValue;
    }
    return array;
}

var colors = shuffle(colorsUnshuffled);
var colorsCursor = 0;

function randColor() {
    return colors[colorsCursor++%colors.length];
}

export default class Logger {
    constructor(origin, env, id, color, prefix) {
        this.env = env;
        this.origin = origin;
        if (color) {
            this.color = color;
        }
        else {
            this.color = randColor();
        }
        if (id) {
            this.id = id;
        }
        else {
            this.id = uid();
        }
        this.prefix = '🤷 %c[Secure Messenger - ' + origin + ' - ' + this.id + ']';
        if (prefix) {
            this.prefix += ' [' + prefix + ']';
        }
    }

    css() {
        return `color: ${this.color};`;
    }

    log(...args) {
        if (this.env === 'development') (console.__log || console.log)(this.prefix, this.css(), ...args);
    }

    warn(...args) {
        if (this.env === 'development') (console.__warn || console.warn)(this.prefix, this.css(), ...args);
    }

    error(...args) {
        (console.__error || console.error)(this.prefix, this.css(), ...args);
    }

    info(...args) {
        if (this.env === 'development') (console.__info || console.info)(this.prefix, this.css(), ...args);
    }

    debug(...args) {
        if (this.env === 'development') (console.__debug || console.debug)(this.prefix, this.css(), ...args);
    }

    addPrefix(prefix) {
        return new Logger(this.origin, this.env, this.id, this.color, prefix);
    }
}