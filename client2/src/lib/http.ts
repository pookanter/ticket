import axios, { type InternalAxiosRequestConfig } from 'axios';
import {
	BehaviorSubject,
	catchError,
	filter,
	firstValueFrom,
	from,
	switchMap,
	take,
	throwError
} from 'rxjs';
import { AuthenService } from './services/authen-service';
import { env } from '$env/dynamic/public';

const refreshTokenState = {
	inProgress: false,
	subject: new BehaviorSubject<unknown>(null)
};

function addAuthenToken(request: InternalAxiosRequestConfig<unknown>) {
	const accessToken = AuthenService.getAccessToken();
	if (!accessToken) {
		return request;
	}

	request.headers['Authorization'] = `Bearer ${accessToken}`;

	return request;
}

export function http() {
	const axiosInstance = axios.create({
		baseURL: env.PUBLIC_API_URL,
		headers: {
			'Content-Type': 'application/json'
		}
	});

	axiosInstance.interceptors.request.use((request) => {
		let baseURL = request.baseURL || '';

		if (request.url?.substring(0, 4) === 'http') {
			baseURL = '';
		} else if (AuthenService.getAccessToken()) {
			request.headers['Authorization'] = `Bearer ${AuthenService.getAccessToken()}`;
		}

		console.log('baseURL', baseURL);

		request.baseURL = baseURL;

		return request;
	});

	axiosInstance.interceptors.response.use(
		(response) => {
			return response;
		},
		(error) => {
			if (axios.isAxiosError(error)) {
				const originalRequest = error.config;

				if (error?.response?.status === 401 && originalRequest) {
					if (refreshTokenState.inProgress) {
						return firstValueFrom(
							refreshTokenState.subject.pipe(
								filter((result) => result !== null),
								take(1),
								switchMap(() => {
									return axiosInstance(addAuthenToken(originalRequest));
								}),
								catchError((err) => {
									return throwError(() => err);
								})
							)
						);
					} else {
						const refreshToken = AuthenService.getRefreshToken();
						if (refreshToken) {
							refreshTokenState.inProgress = true;
							refreshTokenState.subject.next(null);

							return firstValueFrom(
								from(AuthenService.refreshToken(refreshToken)).pipe(
									switchMap((res) => {
										AuthenService.setAuthorization(res.data);
										refreshTokenState.inProgress = false;
										refreshTokenState.subject.next(res);

										return axiosInstance(addAuthenToken(originalRequest));
									}),
									catchError((err) => {
										refreshTokenState.inProgress = false;
										return Promise.reject(err);
									})
								)
							);
						}
					}
				}

				error.message = error.response?.data?.message || error.message;
			}

			return Promise.reject(error);
		}
	);

	return axiosInstance;
}
