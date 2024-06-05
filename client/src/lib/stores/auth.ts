import { goto } from '$app/navigation';
import { AuthenService } from '$lib/services/authen-service';
import { writable, get } from 'svelte/store';

export type AuthState = {
	initializing: boolean;
	user: AuthenService.Me | null | undefined;
};

const authStore = writable<AuthState>({
	initializing: true,
	user: undefined
});

function Use() {
	return authStore.subscribe((value) => {
		if (!value.initializing && value.user == null) {
			console.log('redirect to login');
			goto('/login');
		}
	});
}

export const AuthStore = { ...authStore, Use };
