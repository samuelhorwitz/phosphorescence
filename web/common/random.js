export function shuffle(arr) {
    for (let i = arr.length; i > 0; i--) {
        let randomIndex = Math.floor(Math.random() * i);
        let temporaryValue = arr[i];
        arr[i] = arr[randomIndex];
        arr[randomIndex] = temporaryValue;
    }
    return arr;
}

export function getRandomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}
