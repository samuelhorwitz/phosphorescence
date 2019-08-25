import XRegExp from 'xregexp';

const hashtagMatcher = XRegExp(`(#((?:[\\pL\\pN][\\pM\\u200C\\u200D]*)+(?:[\\p{Pc}\\p{Pd}](?:[\\pL\\pN][\\pM\\u200C\\u200D]*)+)*))`);

export function buildTagMarker(text) {
    let marks = [];
    XRegExp.forEach(text, hashtagMatcher, match => {
        let node = document.createElement('a');
        node.href = `/marketplace/tag/${encodeURIComponent(match[2])}`;
        node.rel = 'tag';
        node.setAttribute('data-bound-html-internal-link', true);
        marks.push({
            index: match.index,
            isOpen: true,
            node
        });
        marks.push({
            index: match.index + match[0].length,
            isOpen: false,
            node
        });
    });
    marks.sort((a, b) => a.index - b.index);
    return marks;
}

export function buildMarker(marks, node) {
    if (!marks || marks.length % 2 !== 0) {
        return null;
    }
    marks.sort((a, b) => a.index - b.index);
    let newMarks = [];
    for (let i = 0; i < marks.length; i += 2) {
        newMarks.push({
            index: marks[i],
            isOpen: true,
            node
        });
        newMarks.push({
            index: marks[i + 1],
            isOpen: false,
            node
        });
    }
    return newMarks;
}

export function combineMarkers(...allMarks) {
    let markMap = new Map();
    let nodeOrder = new WeakMap();
    for (let i = 0; i < allMarks.length; i++) {
        let marks = allMarks[i];
        if (!marks || !marks.length) {
            continue;
        }
        nodeOrder.set(marks[0].node, i);
        for (let j = 0; j < marks.length; j++) {
            let mark = marks[j];
            if (!markMap[mark.index]) {
                markMap[mark.index] = [];
            }
            markMap[mark.index].push(mark);
        }
    }
    let combined = [];
    let openNodes = [];
    for (let key of Object.keys(markMap).sort((a, b) => a - b)) {
        let closers = [];
        let openers = [];
        for (let node of markMap[key]) {
            if (node.isOpen) {
                openers.push(node);
            } else {
                closers.push(node);
            }
        }
        while (closers.length) {
            let found = false;
            for (let i in closers) {
                let node = closers[i];
                if (openNodes[openNodes.length - 1] === node.node) {
                    combined.push(node);
                    openNodes.pop();
                    closers.splice(i, 1);
                    found = true;
                    break;
                }
            }
            if (!found) {
                if (!openNodes.length) {
                    throw new Error('No more open nodes, unbalanced nodes');
                }
                let toReopen = openNodes.pop();
                openers.push({
                    index: parseInt(key, 10),
                    isOpen: true,
                    node: toReopen
                });
                combined.push({
                    index: parseInt(key, 10),
                    isOpen: false,
                    node: toReopen
                });
            }
        }
        for (let node of openers) {
            openNodes.push(node.node);
            combined.push(node);
        }
    }
    return combined;
}

export function getSafeHtml(text, marker) {
    if (!marker || !marker.length) {
        return text;
    }
    let wrapper = document.createElement('div');
    let stack = [];
    let lastIndex = 0;
    for (let mark of marker) {
        let currentEl;
        if (stack.length > 0) {
            currentEl = stack[stack.length - 1];
        } else {
            currentEl = wrapper;
        }
        if (mark.index > lastIndex) {
            let t = document.createTextNode(text.substring(lastIndex, mark.index));
            currentEl.appendChild(t);
        }
        if (mark.isOpen) {
            let el = mark.node.cloneNode();
            currentEl.appendChild(el);
            stack.push(el);
        } else {
            stack.pop();
        }
        lastIndex = mark.index;
    }
    let t = document.createTextNode(text.substring(lastIndex));
    wrapper.appendChild(t);
    return wrapper.innerHTML;
}

export function handleClicks(e) {
    let {target} = e;
    while (target && target.tagName !== 'A') target = target.parentNode;
    if (target && target.getAttribute('data-bound-html-internal-link') === 'true') {
        let {altKey, ctrlKey, metaKey, shiftKey, button, defaultPrevented} = e;
        if (metaKey || altKey || ctrlKey || shiftKey) return;
        if (defaultPrevented) return;
        if (button !== undefined && button !== 0) return;
        if (target && target.getAttribute) {
            let linkTarget = target.getAttribute('target');
            if (/\b_blank\b/i.test(linkTarget)) return;
        }
        let url = new URL(target.href);
        let to = url.pathname;
        if (location.pathname !== to && e.preventDefault) {
            e.preventDefault();
            this.$router.push(to);
        }
    }
}

/*

<em>#Profit-<mark>focused</mark></em>

0 15 em
8 15 mark

0 15 em
0  8 text fragment
8 15 mark


      (0,15)
       /   \
    (0,8) (8,15)


{
    0: {
        15: [em]
    },
    8: {
        15: [mark]
    }
}

[(0, open, em), (8, open, mark), (15, close, mark), (15, close, em)]

create em, set as current parent
create "#Profit-" text node, add to current parent (em)
create mark, add to current parent (em) and set as current parent
create "focused" text node, add to current parent (mark)
pop current parent (dispose mark, back to em)
pop current parent (dispose em, back to root)

<em>foo<mark>ba</em>rrrr</mark>
should rebalance to
<em>foo<mark>ba</mark></em><mark>rrrr</mark>

[(0, open, em), (3, open, mark), (5, close, em), (9, close, mark)]
[(0, open, em), (3, open, mark), (5, close, mark), (5, close, em), (5, open, mark), (9, close, mark)]

create em, set as current parent
create "foo" text node, add to current parent (em)
create mark, add to current parent (em) and set as current parent
create "ba" text node, add to current parent (mark)
pop current parent onto renest stack (mark, back to em)
pop current parent (dispose em, back to root)
pop from renest and set current parent (mark)
create "rrrr" text node, add to current parent (mark)

*/
