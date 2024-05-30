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

export const TicketService = {
	getBoards
};
