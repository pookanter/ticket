<script lang="ts">
	import { goto } from '$app/navigation';
	import SignUpDialogContent from '$lib/components/sign-up-dialog-content/sign-up-dialog-content.svelte';
	// import SignUpContent from '$lib/components/sign-up-content/sign-up-dialog.svelte';
	// import Spinner from '$lib/components/spinner/Spinner.svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	// import DialogFooter from '$lib/components/ui/dialog/dialog-footer.svelte';
	import type { FormInputEvent } from '$lib/components/ui/input';
	import Input from '$lib/components/ui/input/input.svelte';
	import Label from '$lib/components/ui/label/label.svelte';
	import { AuthenService } from '$lib/services/authen-service';
	import authStore from '$lib/stores/auth';
	import { DialogStore } from '$lib/stores/dialog';
	import { AxiosError } from 'axios';
	import { Subject, catchError, from, map, switchMap, throwError } from 'rxjs';
	import { onDestroy, onMount } from 'svelte';
	import type { Unsubscriber } from 'svelte/store';

	let loading = false;
	const alertState = {
		open: false,
		title: '',
		message: ''
	};
	const alertSubject = new Subject();
	const alert$ = alertSubject.subscribe((e: any) => {
		if (e instanceof AxiosError) {
			const { message } = e;
			alertState.open = true;
			alertState.title = 'Error';
			alertState.message = message;
		} else {
			alertState.open = true;
			alertState.title = e.title;
			alertState.message = e.message;
		}
	});

	let unsubscribe: Unsubscriber;
	onMount(() => {
		unsubscribe = authStore.subscribe((state) => {
			console.log('LOGIN MOUNT', state);
			if (state.user) {
				goto('/app');
			}
		});
	});

	onDestroy(() => {
		alert$.unsubscribe();
		unsubscribe();
	});

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

	const signInData = {
		email: '',
		password: ''
	};
	const signInValidation = {
		email: '',
		password: ''
	};
	let signInInvalid = false;

	function handleSignIn(event: Event) {
		event.preventDefault();
		signInInvalid = false;

		if (signInData.email === '') {
			signInValidation.email = 'Email is required';
			signInInvalid = true;
		}

		if (signInData.password === '') {
			signInValidation.password = 'Password is required';
			signInInvalid = true;
		}

		if (signInInvalid) {
			return;
		}

		loading = true;
		from(AuthenService.signIn(signInData))
			.pipe(
				map(({ data }) => {
					AuthenService.setAuthorization(data);
					return data;
				}),
				switchMap(() =>
					from(AuthenService.getMe()).pipe(
						catchError((err) => {
							return throwError(() => err);
						})
					)
				)
			)
			.subscribe({
				next: ({ data: user }) => {
					console.log('success', user);
					loading = false;
					authStore.set({ initializing: false, user });

					goto('/app');
				},
				error: (err) => {
					loading = false;
					alertSubject.next(err);
				}
			});
	}

	// function openSignUp() {
	// 	DialogStore.create(SignUpDialog);
	// }
</script>

<AlertDialog.Root open={alertState.open}>
	<section class="m-auto flex h-[90vh] items-center justify-center">
		<div class="flex h-full flex-shrink-[0.6] flex-col items-center justify-center">
			<div
				class="grid max-w-screen-xl px-4 py-8 mx-auto lg:grid-cols-12 lg:gap-8 lg:py-16 xl:gap-0"
			>
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
									bind:value={signInData.email}
									invalid={signInInvalid}
									errMsg={signInValidation.email}
									placeholder="name@company.com"
								/>
							</div>
							<div>
								<Label class="block mb-2" for="password">Password</Label>
								<Input
									type="password"
									name="password"
									bind:value={signInData.password}
									invalid={signInInvalid}
									errMsg={signInValidation.password}
									placeholder="••••••••"
								/>
							</div>
							<Button class="w-full" on:click={(e) => handleSignIn(e)}>Sign in</Button>
							<p class="text-sm font-light text-gray-500 dark:text-gray-400">
								Don’t have an account yet? <Dialog.Trigger
									on:click={() => {
										DialogStore.create(SignUpDialogContent);
									}}
								>
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

	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>{alertState.title}</AlertDialog.Title>
			<AlertDialog.Description>
				{alertState.message}
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Action>Ok</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
