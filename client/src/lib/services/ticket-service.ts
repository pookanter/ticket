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
	}

	export interface BoardFullDetail {
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

function getBoardById(board_id: number) {
	return http().get<TicketService.BoardFullDetail>(`/${TICKET_SERVICE}/boards/${board_id}`);
}

function createBoard(data: { title: string }) {
	return http().post<TicketService.Board>(`/${TICKET_SERVICE}/boards`, data);
}

function updateBoard({ board_id }: { board_id: number }, data: { title: string }) {
	return http().put<TicketService.Board>(`/${TICKET_SERVICE}/boards/${board_id}`, data);
}

function createStatus({ board_id }: { board_id: number }, data: { title: string }) {
	return http().post<TicketService.Status>(`/${TICKET_SERVICE}/boards/${board_id}/statuses`, data);
}

function updateStatusPartial(
	{ board_id, status_id }: { board_id: number; status_id: number },
	data: {
		title?: string;
		sort_order?: number;
	}
) {
	return http().patch<TicketService.Status>(
		`/${TICKET_SERVICE}/boards/${board_id}/statuses/${status_id}`,
		data
	);
}

function updateStatusesSortOrder(
	{ board_id }: { board_id: number },
	data: { statuses: TicketService.Status[] }
) {
	return http().put<TicketService.Status[]>(
		`/${TICKET_SERVICE}/boards/${board_id}/statuses/sort-orders`,
		data
	);
}

function bulkUpdateTicketOrderInStatuses(
	{ board_id }: { board_id: number },
	data: { statuses: TicketService.Status[] }
) {
	return http().put<TicketService.Ticket[]>(
		`/${TICKET_SERVICE}/boards/${board_id}/statuses/tickets/bulk-reorder`,
		data
	);
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

function updateTicketPartial(
	{ board_id, status_id, ticket_id }: { board_id: number; status_id: number; ticket_id: number },
	data: {
		title?: string;
		description?: string;
		contact?: string;
		status_id?: number;
		sort_order?: number;
	}
) {
	return http().patch<TicketService.Ticket>(
		`/${TICKET_SERVICE}/boards/${board_id}/statuses/${status_id}/tickets/${ticket_id}`,
		data
	);
}

function updateTicketsSortOrder(
	{ board_id, status_id }: { board_id: number; status_id: number },
	data: { tickets: TicketService.Ticket[] }
) {
	return http().put<TicketService.Ticket[]>(
		`/${TICKET_SERVICE}/boards/${board_id}/statuses/${status_id}/tickets/sort-orders`,
		data
	);
}

export const TicketService = {
	getBoards,
	getBoardById,
	createBoard,
	updateBoard,
	createStatus,
	updateStatusPartial,
	updateStatusesSortOrder,
	createTicket,
	updateTicketPartial,
	updateTicketsSortOrder,
	bulkUpdateTicketOrderInStatuses
};
