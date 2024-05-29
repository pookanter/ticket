import type { AuthenService } from '$lib/services/authen-service';
import { writable } from 'svelte/store';

export const auth = writable<AuthenService.Authorization | null>(null);

export const me = writable<AuthenService.Me | null>(null);
