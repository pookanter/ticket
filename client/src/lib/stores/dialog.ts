import * as Dialog from '$lib/components/ui/dialog';
import { type ComponentType } from 'svelte';
import { writable, type Writable } from 'svelte/store';

export type DialogState = {
	initializing: boolean;
	component: ComponentType | null;
};

const dialogStore = writable<DialogState>({
	initializing: true,
	component: null
});

export interface DialogStore extends Writable<DialogState> {
	state: DialogState;
}

export function create(comp: ComponentType) {
	dialogStore.update((store) => {
		store.component = comp;

		return store;
	});
}

export const DialogStore = {
	...dialogStore,
	create
};
