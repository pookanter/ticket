function main() {}

function findTarget(arr: number[], target: number) {
  const sorted = arr.sort((a, b) => a - b);

  for (let i = 0; i < arr.length; i++) {
    for (let j = i; j < arr.length; j++) {
      for (let k = j; k < arr.length; k++) {
        if (sorted[i] + sorted[j] + sorted[k] === target) {
          return [sorted[i] + sorted[j] + sorted[k]];
        }
      }
    }
  }

  return [];
}
