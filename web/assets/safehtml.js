import XRegExp from 'xregexp';
window.XRegExp = XRegExp;

const hashtagMatcher = XRegExp(`(#((?:[\\pL\\pN][\\pM\\u200C\\u200D]*)+(?:[\\p{Pc}\\p{Pd}](?:[\\pL\\pN][\\pM\\u200C\\u200D]*)+)*))`);

export function getSafeHtml(text, marker) {
    if (!marker) {
        return text;
    }
    let el = document.createElement('div');
    let marks = Object.entries(marker);
    if (marks.length === 0) {
        return text;
    }
    marks.sort((a, b) => a[0] - b[0]);
    let lead = document.createTextNode(text.substring(0, marks[0][0]));
    el.appendChild(lead);
    let lastMark;
    for (let [startMark, ends] of marks) {
        if (lastMark) {
            let t = document.createTextNode(text.substring(lastMark, startMark));
            el.appendChild(t);
        }
        let endMarks = Object.entries(ends);
        endMarks.sort((a, b) => a[0] - b[0]);
        let previousEndMark = startMark;
        for (let [endMark, nodes] of endMarks) {
            let previousSubEl;
            for (let node of nodes) {
                let subEl = node.cloneNode();
                if (previousSubEl) {
                    subEl.appendChild(previousSubEl);
                } else {
                    subEl.appendChild(document.createTextNode(text.substring(previousEndMark, endMark)));
                }
                previousSubEl = subEl;
            }
            el.appendChild(previousSubEl);
            previousEndMark = endMark;
            lastMark = endMark;
        }
    }
    let tail = document.createTextNode(text.substring(lastMark));
    el.appendChild(tail);
    return el.innerHTML;
}

export function buildMarker(arrayMark, node, arrayMarks = {}) {
    if (!arrayMark) {
        return arrayMarks;
    }
    for (let i = 0; i < arrayMark.length; i += 2) {
        if (!arrayMarks[arrayMark[i]]) {
            arrayMarks[arrayMark[i]] = {};
        }
        if (!arrayMarks[arrayMark[i]][arrayMark[i + 1]]) {
            arrayMarks[arrayMark[i]][arrayMark[i + 1]] = [];
        }
        arrayMarks[arrayMark[i]][arrayMark[i + 1]].unshift(node);
    }
    return arrayMarks;
}

export function buildTagMarker(text, arrayMarks = {}) {
    let newArrayMarks = [];
    XRegExp.forEach(text, hashtagMatcher, match => {
        newArrayMarks.push(match.index);
        newArrayMarks.push(match.index + match[0].length);
        console.log(match[2]);
    });
    newArrayMarks.sort();
    let node = document.createElement('em');
    return buildMarker(newArrayMarks, node, arrayMarks);
}
