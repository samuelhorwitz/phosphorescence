/* Creative Commons 0 Dedication
 *
 * This work is dedicated under the Creative Commons 0 dedication.
 * To the extent possible under law, the person who associated CC0 with this
 * work has waived all copyright and related or neighboring rights to this work.
 * https://creativecommons.org/publicdomain/zero/1.0/
 *
 * This is contrary to the majority of code in this repository which is licensed
 * under the MIT license with a retained copyright. Only files such as this one
 * which are explicitly licensed differently should be considered licensed under
 * the file-specific license described within. All other files are implicitly
 * licensed under the repository's MIT license.
 */

import {kdTree} from 'kd-tree-javascript';
import {getRandomInt} from '../common/random.js';

export default function wrappedKdTree(points, dimensions, distanceFn) {
    let pointsIndex = {};
    for (let point of points) {
        pointsIndex[point.id] = point;
    }
    let tree = createTree();
    function createTree() {
        return new kdTree(Object.values(pointsIndex), distanceFn, dimensions);
    }
    this.removeById = function removeById(id) {
        delete pointsIndex[id];
        tree = createTree();
    };
    this.nearest = function nearest(k, point) {
        return tree.nearest(point, k).map(([point, distance]) => Object.assign({point, distance}));
    };
    this.getRandomNode = function getRandomNode() {
        let points = Object.values(pointsIndex);
        return {point: points[getRandomInt(0, points.length - 1)]};
    };
    this.getNodesWhere = function getNodesWhere(fn) {
        let points = Object.values(pointsIndex);
        let matches = [];
        for (let point of points) {
            if (fn(point)) {
                matches.push({point});
            }
        }
        return matches;
    };
    this.length = function length() {
        return Object.keys(pointsIndex).length;
    };
    this.getDimensions = function getDimensions() {
        return dimensions;
    };
    this.underlyingTree = tree;
}
