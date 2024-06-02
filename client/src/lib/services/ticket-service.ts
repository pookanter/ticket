import { http } from '$lib/http';

export namespace TicketService {
	export interface Ticket {
		id: number;
		status_id: number;
		title: string;
		description: string;
		contact: string;
		sort_order: number;
		created_at: string;
		updated_at: string;
	}

	export interface Status {
		id: number;
		board_id: number;
		title: string;
		sort_order: number;
		created_at: string;
		updated_at: string;
		tickets: Ticket[];
	}

	export interface Board {
		id: number;
		user_id: number;
		title: string;
		sort_order: number;
		created_at: string;
		updated_at: string;
		statuses: Status[];
	}
}

const TICKET_SERVICE = 'ticket-service' as const;

function getBoards() {
	return http().get<TicketService.Board[]>(`/${TICKET_SERVICE}/boards`);
}

function createBoard(data: { title: string }) {
	return http().post<TicketService.Board>(`/${TICKET_SERVICE}/boards`, data);
}

function createStatus({ board_id }: { board_id: number }, data: { title: string }) {
	return http().post<TicketService.Status>(`/${TICKET_SERVICE}/boards/${board_id}/statuses`, data);
}

function createTicket(
	{ board_id, status_id }: { board_id: number; status_id: number },
	data: { title: string; description: string; contact: string }
) {
	return http().post<TicketService.Ticket>(
		`/${TICKET_SERVICE}/boards/${board_id}/statuses/${status_id}/tickets`,
		data
	);
}

export const TicketService = {
	getBoards,
	createBoard,
	createStatus,
	createTicket
};
