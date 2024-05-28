<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import DialogFooter from '$lib/components/ui/dialog/dialog-footer.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { AuthenService } from '$lib/services/authen-service';
	import { goto } from '$app/navigation';
	import type { FocusProp } from 'bits-ui';
	import Input from '$lib/components/ui/input/input.svelte';
	import Spinner from '$lib/components/spinner/Spinner.svelte';
	import { setAlert } from '$lib/store/app';

	const signupData = {
		email: '',
		name: '',
		lastname: '',
		password: '',
		confirm_password: ''
	};

	let open = false;
	let loading = false;

	function onOpenChange(isOpen: boolean) {
		if (loading) return false;
		open = isOpen;
	}

	async function handleSignup(event: Event) {
		event.preventDefault();
		loading = true;

		try {
			const { data } = await AuthenService.signUp(signupData);

			AuthenService.setAuthorization(data);

			goto('/app');
		} catch ({ error, message }: any) {
			loading = false;
			setAlert({
				title: 'Error',
				message: error ? error.message : message
			});
		}
	}
</script>

<Dialog.Root closeOnOutsideClick={false} {open} {onOpenChange}>
	<Dialog.Trigger>
		<a href={null} class="font-medium text-primary-600 dark:text-primary-500 hover:underline">
			Sign up
		</a>
	</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-[425px]" disableClosing={loading}>
		<Dialog.Header>
			<Dialog.Title>Sign up</Dialog.Title>
		</Dialog.Header>
		<div class="p-4 md:p-5">
			<form class="space-y-4 md:space-y-6" action="#">
				<div>
					<label for="email" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
						>Email</label
					>
					<Input
						type="text"
						placeholder="name@company.com"
						class="block w-full"
						disabled={loading}
						bind:value={signupData.email}
					/>
				</div>
				<div>
					<label for="email" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
						>Name</label
					>
					<Input
						type="text"
						placeholder="John"
						class="block w-full"
						disabled={loading}
						bind:value={signupData.name}
					/>
				</div>
				<div>
					<label for="email" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
						>Last name</label
					>
					<Input
						type="text"
						placeholder="Doe"
						class="block w-full"
						disabled={loading}
						bind:value={signupData.lastname}
					/>
				</div>
				<div>
					<label for="password" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
						>Password</label
					>
					<Input
						type="password"
						placeholder="••••••••"
						class="block w-full"
						disabled={loading}
						bind:value={signupData.password}
					/>
				</div>
				<div>
					<label for="password" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
						>Confirm password</label
					>
					<Input
						type="confirm_password"
						placeholder="••••••••"
						class="block w-full"
						disabled={loading}
						bind:value={signupData.confirm_password}
					/>
				</div>
			</form>
		</div>
		<DialogFooter>
			<Button class="w-full" on:click={(e) => handleSignup(e)}>
				{#if loading}
					<Spinner class="w-4 h-4 mr-2" /> Signing up
				{:else}
					Sign up
				{/if}
			</Button>
		</DialogFooter>
	</Dialog.Content>
</Dialog.Root>
