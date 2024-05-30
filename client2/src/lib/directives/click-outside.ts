import type { ActionReturn } from 'svelte/action';

interface Attributes {
	'on:clickoutside'?: (e: CustomEvent<void>) => void;
}

type Callback = () => unknown;

export function clickOutsideAction(
	node: HTMLElement,
	callback?: Callback
): ActionReturn<{}, Attributes> {
	const handleClick = (event: MouseEvent) => {
		if (node && !node.contains(event.target as Node) && !event.defaultPrevented) {
			node.dispatchEvent(new CustomEvent('clickoutside', node as any));
		}
	};

	document.addEventListener('click', handleClick, true);

	return {
		destroy() {
			document.removeEventListener('click', handleClick, true);
		}
	};
}
