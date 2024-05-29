<script lang="ts">
	import Spinner from '$lib/components/spinner/Spinner.svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import DialogFooter from '$lib/components/ui/dialog/dialog-footer.svelte';
	import Input from '$lib/components/ui/input/input.svelte';
	import Label from '$lib/components/ui/label/label.svelte';

	let loading = false;
	const signInData = {
		email: '',
		password: ''
	};
	const signUpData = {
		email: '',
		name: '',
		lastname: '',
		password: '',
		confirm_password: ''
	};
	const signUpValidation = {
		email: '',
		name: '',
		lastname: '',
		password: '',
		confirm_password: ''
	};
	const alertData = {
		open: false,
		title: '',
		message: ''
	};

	function handleSignIn(event: Event) {
		event.preventDefault();
	}
	function handleSignUp(event: Event) {
		event.preventDefault();

		if (signUpData.name === '') {
			signUpValidation.name = 'Name is required';
			return;
		}

		signUpValidation.name = '';

		if (signUpData.lastname === '') {
			signUpValidation.lastname = 'Lastname is required';
			return;
		}

		signUpValidation.lastname = '';

		if (signUpData.password === '') {
			signUpValidation.password = 'Password is required';
			return;
		}

		signUpValidation.password = '';

		if (signUpData.confirm_password === '') {
			signUpValidation.confirm_password = 'Confirm password is required';
			return;
		}

		if (signUpData.password !== signUpData.confirm_password) {
			signUpValidation.confirm_password = 'Password and confirm password must be the same';
			return;
		}

		signUpValidation.confirm_password = '';
	}

	let open = false;
	function onOpenChange(isOpen: boolean) {
		if (loading) return false;
		open = isOpen;

		if (!isOpen) {
			let key: keyof typeof signUpData;
			for (key in signUpData) {
				signUpData[key] = '';
				signUpValidation[key] = '';
			}
		}
	}
</script>

<AlertDialog.Root>
	<Dialog.Root {open} {onOpenChange}>
		<section class="m-auto flex h-[90vh] items-center justify-center">
			<div class="flex h-full flex-shrink-[0.6] flex-col items-center justify-center">
				<div
					class="grid max-w-screen-xl px-4 py-8 mx-auto lg:grid-cols-12 lg:gap-8 lg:py-16 xl:gap-0"
				>
					<div class="mr-auto place-self-center lg:col-span-7">
						<h1
							class="max-w-2xl mb-4 text-4xl font-extrabold leading-none tracking-tight md:text-5xl xl:text-6xl dark:text-white"
						>
							Streamline Your Support Tickets with Ease
						</h1>
						<p
							class="max-w-2xl mb-6 font-light text-gray-500 md:text-lg lg:mb-8 lg:text-xl dark:text-gray-400"
						>
							Transform the way you manage customer support. Our Kanban-style ticketing system
							simplifies workflows, boosts productivity, and ensures nothing slips through the
							cracks.
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
										bind:value={signInData.email}
										placeholder="name@company.com"
									/>
								</div>
								<div>
									<Label class="block mb-2" for="password">Password</Label>
									<Input
										type="password"
										name="password"
										bind:value={signInData.password}
										placeholder="••••••••"
									/>
								</div>
								<Button class="w-full" on:click={(e) => handleSignIn(e)}>Sign in</Button>
								<p class="text-sm font-light text-gray-500 dark:text-gray-400">
									Don’t have an account yet? <Dialog.Trigger>
										<a
											href={null}
											class="font-medium text-primary-600 dark:text-primary-500 hover:underline"
										>
											Sign up
										</a>
									</Dialog.Trigger>
								</p>
							</form>
						</div>
					</div>
				</div>
			</div>
		</section>

		<Dialog.Content class="sm:max-w-[425px]" disableClosing={loading}>
			<Dialog.Header>
				<Dialog.Title>Sign up</Dialog.Title>
			</Dialog.Header>
			<form class="space-y-4 md:space-y-6" action="#">
				<div>
					<Label class="block mb-2" for="terms">Your email</Label>
					<Input
						type="email"
						name="email"
						bind:value={signUpData.email}
						disabled={loading}
						placeholder="name@company.com"
					/>
				</div>
				<div>
					<Label class="block mb-2" for="terms">Name</Label>
					<Input
						type="name"
						name="name"
						bind:value={signUpData.name}
						disabled={loading}
						placeholder="John"
						errMsg={signUpValidation.name}
					/>
				</div>
				<div>
					<Label class="block mb-2" for="terms">Lastname</Label>
					<Input
						type="name"
						name="name"
						bind:value={signUpData.lastname}
						disabled={loading}
						placeholder="Smith"
						errMsg={signUpValidation.lastname}
					/>
				</div>
				<div>
					<Label class="block mb-2" for="password">Password</Label>
					<Input
						type="password"
						name="password"
						bind:value={signUpData.password}
						disabled={loading}
						placeholder="••••••••"
						errMsg={signUpValidation.password}
					/>
				</div>
				<div>
					<Label class="block mb-2" for="confirm_password">Confirm password</Label>
					<Input
						type="confirm_password"
						name="confirm_password"
						bind:value={signUpData.confirm_password}
						disabled={loading}
						placeholder="••••••••"
						errMsg={signUpValidation.confirm_password}
					/>
				</div>
			</form>
			<DialogFooter>
				<Button class="w-full" on:click={(e) => handleSignUp(e)}>
					{#if loading}
						<Spinner class="w-4 h-4 mr-2" /> Signing up
					{:else}
						Sign up
					{/if}
				</Button>
			</DialogFooter>
		</Dialog.Content>
	</Dialog.Root>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>{alertData.title}</AlertDialog.Title>
			<AlertDialog.Description>
				{alertData.message}
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Action>Ok</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

<!-- <div class='flex items-center justify-center h-screen space-x-2 bg-white dark:invert'>
  <span class='sr-only'>Loading...</span>
   <div class='h-8 w-8 bg-black rounded-full animate-bounce [animation-delay:-0.3s]'></div>
 <div class='h-8 w-8 bg-black rounded-full animate-bounce [animation-delay:-0.15s]'></div>
 <div class='w-8 h-8 bg-black rounded-full animate-bounce'></div>
</div> -->
