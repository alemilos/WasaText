/**
 * Create a "hacker" style animation to characters on the input element
 * @param {*} el the element with text
 */
export function hackerText(el) {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";
	let iterations = 0;
	let text = el.innerText;
	const interval = setInterval(() => {
		el.innerText = text
			.split("")
			.map((letter, index) => {
				function upperOrLower(current, original) {
					if (original.toUpperCase() === original) {
						return current.toUpperCase();
					}
					return current.toLowerCase();
				}

				if (index < iterations || letter === " ") return letter;

				return upperOrLower(
					letters[Math.floor(Math.random() * 26)],
					letter
				);
			})
			.join("");

		if (iterations >= text.length) clearInterval(interval);
		iterations += 1 / 3;
	}, 30);
}
