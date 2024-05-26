<script lang="ts">
	import SignUp from './SignUp.svelte';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import DialogDescription from '$lib/components/ui/dialog/dialog-description.svelte';
	import Header from './Header.svelte';
	import DialogFooter from '$lib/components/ui/dialog/dialog-footer.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
</script>

<div class="p-6 space-y-4 bg-white rounded-lg sm:p-8 md:space-y-6">
	<h1
		class="text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white"
	>
		Sign in
	</h1>
	<form class="space-y-4 md:space-y-6 lg:min-w-80" action="#">
		<div>
			<label for="email" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
				>Your email</label
			>
			<input
				type="email"
				name="email"
				id="email"
				class="focus:ring-primary-600 focus:border-primary-600 block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-gray-900 sm:text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400 dark:focus:border-blue-500 dark:focus:ring-blue-500"
				placeholder="name@company.com"
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
				class="focus:ring-primary-600 focus:border-primary-600 block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-gray-900 sm:text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400 dark:focus:border-blue-500 dark:focus:ring-blue-500"
			/>
		</div>
		<Button class="w-full">Sign in</Button>
		<p class="text-sm font-light text-gray-500 dark:text-gray-400">
			Don’t have an account yet?
			<Dialog.Root>
				<Dialog.Trigger
					><a href={null} class="font-medium text-primary-600 dark:text-primary-500 hover:underline"
						>Sign up</a
					></Dialog.Trigger
				>

				<Dialog.Content class="sm:max-w-[425px]">
					<Dialog.Header>
						<Dialog.Title>Sign up</Dialog.Title>
					</Dialog.Header>
					<div class="p-4 md:p-5">
						<form class="space-y-4 md:space-y-6" action="#">
							<div>
								<label
									for="email"
									class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Email</label
								>
								<input
									type="email"
									name="email"
									id="email"
									class="block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-500 dark:bg-gray-600 dark:text-white dark:placeholder-gray-400"
									placeholder="name@company.com"
									required
								/>
							</div>
							<div>
								<label
									for="email"
									class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Name</label
								>
								<input
									type="text"
									name="name"
									id="name"
									class="block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-500 dark:bg-gray-600 dark:text-white dark:placeholder-gray-400"
									placeholder="John"
									required
								/>
							</div>
							<div>
								<label
									for="email"
									class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
									>Last name</label
								>
								<input
									type="text"
									name="lastname"
									id="lastname"
									class="block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-500 dark:bg-gray-600 dark:text-white dark:placeholder-gray-400"
									placeholder="Doe"
									required
								/>
							</div>
							<div>
								<label
									for="password"
									class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
									>Password</label
								>
								<input
									type="password"
									name="password"
									id="password"
									placeholder="••••••••"
									class="block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-500 dark:bg-gray-600 dark:text-white dark:placeholder-gray-400"
									required
								/>
							</div>
							<div>
								<label
									for="password"
									class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
									>Confirm password</label
								>
								<input
									type="password"
									name="password"
									id="password"
									placeholder="••••••••"
									class="block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-500 dark:bg-gray-600 dark:text-white dark:placeholder-gray-400"
									required
								/>
							</div>
						</form>
					</div>
					<DialogFooter>
						<Button class="w-full">Sign up</Button>
					</DialogFooter>
				</Dialog.Content>
			</Dialog.Root>
		</p>
	</form>
</div>
