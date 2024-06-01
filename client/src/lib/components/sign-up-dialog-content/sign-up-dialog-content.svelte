<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import Input from '$lib/components/ui/input/input.svelte';
	import Label from '$lib/components/ui/label/label.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { AuthenService } from '$lib/services/authen-service';
	import Spinner from '$lib/components/spinner/Spinner.svelte';
	import type { FormInputEvent } from '../ui/input';
	const data = {
		email: '',
		name: '',
		lastname: '',
		password: '',
		confirm_password: ''
	};
	const error = {
		email: '',
		name: '',
		lastname: '',
		password: '',
		confirm_password: ''
	};
	let invalid = false;
	let loading = false;

	// const t: Dialog.Content = Dialog.Content

	function clear() {
		let key: keyof typeof error;
		for (key in error) {
			error[key] = '';
		}
	}

	async function handleSignUp(event: Event) {
		event.preventDefault();
		clear();
		invalid = false;

		if (data.email === '') {
			error.email = 'Email is required';
			invalid = true;
		}

		if (data.name === '') {
			error.name = 'Name is required';
			invalid = true;
		}
		if (data.lastname === '') {
			error.lastname = 'Lastname is required';
			invalid = true;
		}

		if (data.password === '') {
			error.password = 'Password is required';
			invalid = true;
		}

		if (data.password.length < 8) {
			error.password = 'Password must be at least 8 characters';
			invalid = true;
		}

		if (data.confirm_password === '') {
			error.confirm_password = 'Confirm password is required';
			invalid = true;
		}

		if (data.password !== data.confirm_password) {
			error.confirm_password = 'Password and confirm password must be the same';
			invalid = true;
		}

		if (invalid) {
			return;
		}

		loading = true;
		try {
			await AuthenService.signUp(data);
		} catch (error) {}
		loading = false;
	}

	function disableWhitespace(e: FormInputEvent<KeyboardEvent>) {
		if (e?.code === 'Space') {
			e.preventDefault();
		}
	}

	function noWhitespaceFuncs<Key extends PropertyKey>(src: { [K in Key]: string }, key: Key) {
		return (e: any) => {
			src[key] = (src[key] as string).replaceAll(' ', '');
		};
	}
</script>

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
				bind:value={data.email}
				on:keydown={disableWhitespace}
				on:input={noWhitespaceFuncs(data, 'email')}
				disabled={loading}
				placeholder="name@company.com"
				{invalid}
				errMsg={error.email}
			/>
		</div>
		<div>
			<Label class="block mb-2" for="terms">Name</Label>
			<Input
				type="name"
				name="name"
				bind:value={data.name}
				on:keydown={disableWhitespace}
				on:change={noWhitespaceFuncs(data, 'name')}
				disabled={loading}
				placeholder="John"
				{invalid}
				errMsg={error.name}
			/>
		</div>
		<div>
			<Label class="block mb-2" for="terms">Lastname</Label>
			<Input
				type="lastname"
				name="lastname"
				bind:value={data.lastname}
				on:keydown={disableWhitespace}
				on:input={noWhitespaceFuncs(data, 'lastname')}
				disabled={loading}
				placeholder="Smith"
				{invalid}
				errMsg={error.lastname}
			/>
		</div>
		<div>
			<Label class="block mb-2" for="password">Password</Label>
			<Input
				type="password"
				name="password"
				bind:value={data.password}
				disabled={loading}
				placeholder="••••••••"
				{invalid}
				errMsg={error.password}
			/>
		</div>
		<div>
			<Label class="block mb-2" for="confirm_password">Confirm password</Label>
			<Input
				type="password"
				name="confirm_password"
				bind:value={data.confirm_password}
				disabled={loading}
				placeholder="••••••••"
				{invalid}
				errMsg={error.confirm_password}
			/>
		</div>
	</form>
	<Dialog.Footer>
		<Button class="w-full" on:click={(e) => handleSignUp(e)}>
			{#if loading}
				<Spinner class="w-4 h-4 mr-2" /> Signing up
			{:else}
				Sign up
			{/if}
		</Button>
	</Dialog.Footer>
</Dialog.Content>
