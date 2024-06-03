<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import '../app.css';
	import authStore, { type AuthStore } from '$lib/stores/auth';
	import { AuthenService } from '$lib/services/authen-service';
	import { goto } from '$app/navigation';
	import type { Unsubscriber } from 'svelte/store';
	import * as Dialog from '$lib/components/ui/dialog';
	import { DialogStore } from '$lib/stores/dialog';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { AlertStore } from '$lib/stores/alert';
	import { ModeWatcher, toggleMode } from 'mode-watcher';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Moon, Sun } from 'lucide-svelte';

	let unsubscribes: Unsubscriber[] = [];
	let alertState = AlertStore.defaultState();
	let dialogState = DialogStore.defaultState();
	onMount(() => {
		unsubscribes.push(
			authStore.subscribe(async (state) => {
				if (state.initializing) {
					console.log('state.initializing...');
					const token = AuthenService.getAuthorization();
					const storeValue: AuthStore = {
						initializing: false,
						user: null
					};
					if (token) {
						try {
							const { data: user } = await AuthenService.getMe();

							storeValue.user = user;
						} catch (error) {
							console.error(error);
						}
					}

					authStore.set(storeValue);

					// setTimeout(() => {
					// 	if (storeValue.user) {
					// 		goto('/');
					// 	} else {
					// 		goto('/login');
					// 	}
					// }, 0);
				}
			})
		);

		unsubscribes.push(
			DialogStore.subscribe((state) => {
				if (state.initializing) {
					DialogStore.update((store) => {
						store.initializing = false;

						return store;
					});

					return;
				}

				dialogState = state;
			})
		);

		unsubscribes.push(
			AlertStore.subscribe((state) => {
				if (state.initializing) {
					AlertStore.update((store) => {
						store.initializing = false;

						return store;
					});

					return;
				}

				alertState = state;
			})
		);
	});

	onDestroy(() => {
		unsubscribes.forEach((unsubscribe) => {
			unsubscribe();
		});
	});
</script>

<ModeWatcher />
<div class="relative flex flex-col justify-between min-h-screen bg-background">
	<header
		class="box-content h-[var(--header-height)] z-50 w-full border-b border-border/40 bg-background/95"
	>
		<div class="flex items-center justify-end h-full px-3">
			<Button on:click={toggleMode} variant="outline" size="icon">
				<Sun
					class="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0"
				/>
				<Moon
					class="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100"
				/>
				<span class="sr-only">Toggle theme</span>
			</Button>
		</div>
	</header>
	<main class="flex items-stretch flex-1">
		<AlertDialog.Root
			open={alertState.open}
			onOpenChange={alertState.onOpenChange}
			closeOnOutsideClick={false}
		>
			<AlertDialog.Content class="z-[60]">
				<AlertDialog.Header>
					{#if alertState.title}
						<AlertDialog.Title>{alertState.title}</AlertDialog.Title>
					{/if}
					{#if alertState.message}
						<AlertDialog.Description>{alertState.message}</AlertDialog.Description>
					{/if}
				</AlertDialog.Header>
				<AlertDialog.Footer>
					{#each alertState.buttons as button}
						{#if button.type === 'cancel'}
							<AlertDialog.Cancel on:click={button.onClick}>{button.text}</AlertDialog.Cancel>
						{:else}
							<AlertDialog.Action on:click={button.onClick}>{button.text}</AlertDialog.Action>
						{/if}
					{/each}
				</AlertDialog.Footer>
			</AlertDialog.Content>
			<Dialog.Root
				open={dialogState.open}
				onOpenChange={dialogState.onOpenChange}
				closeOnOutsideClick={false}
			>
				<slot />
				{#if dialogState.component}
					<svelte:component this={dialogState.component} {...dialogState.params} />
				{/if}
			</Dialog.Root>
		</AlertDialog.Root>
	</main>
	<footer class="box-content w-full bg-background/95 h-[var(--footer-height)]">
		<div class="w-full max-w-screen-xl p-4 mx-auto md:flex md:items-center md:justify-between">
			<span class="text-sm text-gray-500 sm:text-center dark:text-gray-400"
				>Made by <b>Sippakorn Phuakpong</b> with SvelteKit and TailwindCSS with back-end write in Go.
			</span>
			<ul
				class="flex flex-wrap items-center mt-3 text-sm font-medium text-gray-500 sm:mt-0 dark:text-gray-400"
			>
				<li>
					<a
						href="https://github.com/sippakorn-phuakpong"
						target="_blank"
						class="flex items-end me-4 hover:underline md:me-6"
					>
						<svg
							class="w-6 h-6 text-gray-800 dark:text-white"
							aria-hidden="true"
							xmlns="http://www.w3.org/2000/svg"
							width="24"
							height="24"
							fill="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								fill-rule="evenodd"
								d="M12.006 2a9.847 9.847 0 0 0-6.484 2.44 10.32 10.32 0 0 0-3.393 6.17 10.48 10.48 0 0 0 1.317 6.955 10.045 10.045 0 0 0 5.4 4.418c.504.095.683-.223.683-.494 0-.245-.01-1.052-.014-1.908-2.78.62-3.366-1.21-3.366-1.21a2.711 2.711 0 0 0-1.11-1.5c-.907-.637.07-.621.07-.621.317.044.62.163.885.346.266.183.487.426.647.71.135.253.318.476.538.655a2.079 2.079 0 0 0 2.37.196c.045-.52.27-1.006.635-1.37-2.219-.259-4.554-1.138-4.554-5.07a4.022 4.022 0 0 1 1.031-2.75 3.77 3.77 0 0 1 .096-2.713s.839-.275 2.749 1.05a9.26 9.26 0 0 1 5.004 0c1.906-1.325 2.74-1.05 2.74-1.05.37.858.406 1.828.101 2.713a4.017 4.017 0 0 1 1.029 2.75c0 3.939-2.339 4.805-4.564 5.058a2.471 2.471 0 0 1 .679 1.897c0 1.372-.012 2.477-.012 2.814 0 .272.18.592.687.492a10.05 10.05 0 0 0 5.388-4.421 10.473 10.473 0 0 0 1.313-6.948 10.32 10.32 0 0 0-3.39-6.165A9.847 9.847 0 0 0 12.007 2Z"
								clip-rule="evenodd"
							/>
						</svg>

						<span class="ml-2">Github</span>
					</a>
				</li>
				<li>
					<a
						href="https://www.linkedin.com/in/sippakorn-phuakpong-758031295/"
						target="_blank"
						class="flex items-end hover:underline"
					>
						<svg
							class="w-6 h-6 text-gray-800 dark:text-white"
							aria-hidden="true"
							xmlns="http://www.w3.org/2000/svg"
							width="24"
							height="24"
							fill="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								fill-rule="evenodd"
								d="M12.51 8.796v1.697a3.738 3.738 0 0 1 3.288-1.684c3.455 0 4.202 2.16 4.202 4.97V19.5h-3.2v-5.072c0-1.21-.244-2.766-2.128-2.766-1.827 0-2.139 1.317-2.139 2.676V19.5h-3.19V8.796h3.168ZM7.2 6.106a1.61 1.61 0 0 1-.988 1.483 1.595 1.595 0 0 1-1.743-.348A1.607 1.607 0 0 1 5.6 4.5a1.601 1.601 0 0 1 1.6 1.606Z"
								clip-rule="evenodd"
							/>
							<path d="M7.2 8.809H4V19.5h3.2V8.809Z" />
						</svg>
						<span>Linkedin</span>
					</a>
				</li>
			</ul>
		</div>
	</footer>
</div>
