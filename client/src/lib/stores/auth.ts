import type { AuthenService } from '$lib/services/authen-service';
import { writable } from 'svelte/store';

export type AuthStore = {
	initializing: boolean;
	user: AuthenService.Me | null | undefined;
};

const authStore = writable<AuthStore>({
	initializing: true,
	user: undefined
});

export default authStore;
