import { type ComponentType } from 'svelte';
import { writable, type Writable } from 'svelte/store';

type Data = any;

type Buttons = Array<{ type: 'cancel' | 'action'; text: string; onClick: () => void }>;

export type AlertState<T> = {
	initializing: boolean;
	title: string;
	message: string;
	open: boolean;
	forceClose: boolean;
	data?: T;
	buttons: Buttons;
	onOpenChange: (isOpen: boolean) => void;
	onClose?: (d: T) => boolean;
};

function onOpenChange(isOpen: boolean) {
	console.log('onOpenChange', isOpen);
	alertStore.update((store) => {
		if (!isOpen && !store.forceClose && store.onClose) {
			const closeable = store.onClose(store.data);
			store.open = closeable == false;
		} else {
			store.open = isOpen;
		}

		return store;
	});
}

function defaultState(): AlertState<Data> {
	return {
		initializing: true,
		title: '',
		message: '',
		open: false,
		forceClose: false,
		buttons: [
			{
				type: 'action',
				text: 'ok',
				onClick: () => {}
			}
		],
		onOpenChange: onOpenChange
	};
}

const alertStore = writable<AlertState<Data>>(defaultState());

export interface AlertStore<T> extends Writable<AlertState<T>> {
	state: AlertState<T>;
}

type Param = {
	title?: string;
	message?: string;
	buttons?: Buttons;
};

function create({ title, message, buttons }: Param) {
	alertStore.update((store) => {
		store.title = title || '';
		store.message = message || '';
		store.buttons = buttons || store.buttons;
		store.open = true;

		return store;
	});
}

export const AlertStore = {
	...alertStore,
	defaultState,
	create,
	close
};
