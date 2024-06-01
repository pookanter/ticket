<script lang="ts">
	import { goto } from '$app/navigation';
	import SignUpDialogContent from '$lib/components/sign-up-dialog-content/sign-up-dialog-content.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import Input from '$lib/components/ui/input/input.svelte';
	import Label from '$lib/components/ui/label/label.svelte';
	import { AuthenService } from '$lib/services/authen-service';
	import { AlertStore } from '$lib/stores/alert';
	import authStore from '$lib/stores/auth';
	import { DialogStore } from '$lib/stores/dialog';
	import { onDestroy, onMount } from 'svelte';
	import type { Unsubscriber } from 'svelte/store';

	let loading = false;

	let unsubscribe: Unsubscriber;
	onMount(() => {
		unsubscribe = authStore.subscribe((state) => {
			console.log('LOGIN MOUNT', state);
			if (state.user) {
				goto('/');
			}
		});
	});

	onDestroy(() => {
		unsubscribe();
	});

	const data = {
		email: '',
		password: ''
	};
	const error = {
		email: '',
		password: ''
	};
	let invalid = false;

	async function handleSignIn(event: Event) {
		event.preventDefault();
		invalid = false;

		if (data.email === '') {
			error.email = 'Email is required';
			invalid = true;
		}

		if (data.password === '') {
			error.password = 'Password is required';
			invalid = true;
		}

		if (invalid) {
			return;
		}

		loading = true;
		try {
			const { data: token } = await AuthenService.signIn(data);

			AuthenService.setAuthorization(token);

			const { data: user } = await AuthenService.getMe();

			authStore.set({ initializing: false, user });
			goto('/');
		} catch ({ error, message }: any) {
			AlertStore.create({
				title: 'Error',
				message: message
			});
		}
		loading = false;
	}
</script>

<section class="m-auto flex h-[90vh] items-center justify-center">
	<div class="flex h-full flex-shrink-[0.6] flex-col items-center justify-center">
		<div class="grid max-w-screen-xl px-4 py-8 mx-auto lg:grid-cols-12 lg:gap-8 lg:py-16 xl:gap-0">
			<div class="mr-auto place-self-center lg:col-span-7">
				<h1
					class="max-w-2xl mb-4 text-4xl font-extrabold leading-none tracking-tight md:text-5xl xl:text-6xl dark:text-white cursor-pointer bg-[var(--theme-color)]"
				>
					Streamline Your Support Tickets with Ease
				</h1>
				<p
					class="max-w-2xl mb-6 font-light text-gray-500 md:text-lg lg:mb-8 lg:text-xl dark:text-gray-400 bg-[var(--theme-color)]"
				>
					Transform the way you manage customer support. Our Kanban-style ticketing system
					simplifies workflows, boosts productivity, and ensures nothing slips through the cracks.
				</p>
			</div>
			<div class="flex justify-center mt-8 lg:col-span-5 lg:mt-0 lg:justify-end">
				<div class="p-6 space-y-4 bg-white rounded-lg sm:p-8 md:space-y-6">
					<h1
						class="text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white"
					>
						Sign in
					</h1>
					<form class="space-y-4 md:space-y-6 lg:min-w-80" action="#">
						<div>
							<Label class="block mb-2" for="terms">Your email</Label>
							<Input
								type="email"
								name="email"
								bind:value={data.email}
								{invalid}
								errMsg={error.email}
								placeholder="name@company.com"
							/>
						</div>
						<div>
							<Label class="block mb-2" for="password">Password</Label>
							<Input
								type="password"
								name="password"
								bind:value={data.password}
								{invalid}
								errMsg={error.password}
								placeholder="••••••••"
							/>
						</div>
						<Button class="w-full" on:click={(e) => handleSignIn(e)}>Sign in</Button>
						<p class="text-sm font-light text-gray-500 dark:text-gray-400">
							Don’t have an account yet?
							<a
								href={null}
								class="font-medium cursor-pointer text-primary-600 dark:text-primary-500 hover:underline"
								on:click={() => {
									DialogStore.create(SignUpDialogContent);
								}}
							>
								Sign up
							</a>
						</p>
					</form>
				</div>
			</div>
		</div>
	</div>
</section>
