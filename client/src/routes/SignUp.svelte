<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import DialogFooter from '$lib/components/ui/dialog/dialog-footer.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { AuthenService } from '$lib/services/authen-service';
	import { goto } from '$app/navigation';

	const signupData = {
		email: '',
		name: '',
		lastname: '',
		password: '',
		confirm_password: ''
	};

	async function handleSignup(event: Event) {
		event.preventDefault();
		try {
			const { data } = await AuthenService.signUp(signupData);

			AuthenService.setAuthorization(data);

			goto('/app');
		} catch (error) {
			console.error(error);

			// Dialog.show({
			//   title: 'Error',
			//   description: error.message
			// });
		}
	}
</script>

<div class="p-4 md:p-5">
	<form class="space-y-4 md:space-y-6" action="#">
		<div>
			<label for="email" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
				>Email</label
			>
			<input
				type="email"
				name="email"
				id="email"
				bind:value={signupData.email}
				class="block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-500 dark:bg-gray-600 dark:text-white dark:placeholder-gray-400"
				placeholder="name@company.com"
				required
			/>
		</div>
		<div>
			<label for="email" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
				>Name</label
			>
			<input
				type="text"
				name="name"
				id="name"
				bind:value={signupData.name}
				class="block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-500 dark:bg-gray-600 dark:text-white dark:placeholder-gray-400"
				placeholder="John"
				required
			/>
		</div>
		<div>
			<label for="email" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
				>Last name</label
			>
			<input
				type="text"
				name="lastname"
				id="lastname"
				bind:value={signupData.lastname}
				class="block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-500 dark:bg-gray-600 dark:text-white dark:placeholder-gray-400"
				placeholder="Doe"
				required
			/>
		</div>
		<div>
			<label for="password" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
				>Password</label
			>
			<input
				type="password"
				name="password"
				id="password"
				placeholder="••••••••"
				bind:value={signupData.password}
				class="block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-500 dark:bg-gray-600 dark:text-white dark:placeholder-gray-400"
				required
			/>
		</div>
		<div>
			<label for="password" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
				>Confirm password</label
			>
			<input
				type="password"
				name="password"
				id="password"
				placeholder="••••••••"
				bind:value={signupData.confirm_password}
				class="block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-500 dark:bg-gray-600 dark:text-white dark:placeholder-gray-400"
				required
			/>
		</div>
	</form>
</div>
