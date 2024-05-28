import { atom } from 'nanostores';

export const $alert = atom<{
	title: string;
	message: string;
} | null>(null);

export function setAlert(msg: { title: string; message: string }) {
	$alert.set(msg);
}
