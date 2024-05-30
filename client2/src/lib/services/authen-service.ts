import { http } from '$lib/http';

export namespace AuthenService {
	export interface Authorization {
		access_token: string;
		refresh_token: string;
	}

	export interface SignInRequest {
		email: string;
		password: string;
	}

	export interface SignUpRequest {
		name: string;
		lastname: string;
		email: string;
		password: string;
	}

	export interface Me {
		id: string;
		name: string;
		lastname: string;
		email: string;
		creted_at: string;
	}
}

const AUTHEN_SERVICE = 'authen-service' as const;

function getAuthorization() {
	const token = localStorage.getItem('x-authorization');

	if (token) {
		return JSON.parse(token) as AuthenService.Authorization;
	}

	return null;
}

function setAuthorization(token: AuthenService.Authorization) {
	localStorage.setItem('x-authorization', JSON.stringify(token));
}

function removeAuthorization() {
	localStorage.removeItem('x-authorization');
}

function getAccessToken() {
	const token = getAuthorization();

	return token?.access_token;
}

function getRefreshToken() {
	const token = getAuthorization();

	return token?.refresh_token;
}

function refreshToken(refresh_token: string) {
	return http().post<AuthenService.Authorization>(`/${AUTHEN_SERVICE}/refresh-token`, {
		refresh_token
	});
}

function signIn(data: AuthenService.SignInRequest) {
	return http().post<AuthenService.Authorization>(`/${AUTHEN_SERVICE}/sign-in`, data);
}

function signUp(data: AuthenService.SignUpRequest) {
	return http().post(`/${AUTHEN_SERVICE}/sign-up`, data);
}

function getMe() {
	return http().get<AuthenService.Me>(`/${AUTHEN_SERVICE}/users/me`);
}

export const AuthenService = {
	AUTHEN_SERVICE,
	getAuthorization,
	setAuthorization,
	removeAuthorization,
	getAccessToken,
	getRefreshToken,
	refreshToken,
	signIn,
	signUp,
	getMe
};
