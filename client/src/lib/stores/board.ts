import { TicketService } from '$lib/services/ticket-service';
import { cloneDeep } from 'lodash';
import { from } from 'rxjs';
import { writable, get } from 'svelte/store';

export type BoardState = {
	selected?: TicketService.BoardFullDetail;
	boards: TicketService.Board[];
};

function defaultState(): BoardState {
	return {
		selected: undefined,
		boards: []
	};
}

const boardStore = writable<BoardState>(defaultState());

function selectBoard(board: TicketService.BoardFullDetail) {
	boardStore.update((state) => {
		const selected = cloneDeep(board);
		selected.statuses = selected.statuses.sort((a, b) => a.sort_order - b.sort_order);

		for (let i = 0; i < selected.statuses.length; i++) {
			selected.statuses[i].tickets = selected.statuses[i].tickets.sort(
				(a, b) => a.sort_order - b.sort_order
			);
		}

		state.selected = selected;

		return state;
	});
}

function addBoard({ board }: { board: TicketService.Board }) {
	boardStore.update((state) => {
		const boards = cloneDeep(state.boards);
		boards.push(board);

		state.boards = boards;

		return state;
	});
}

function addStatus({ status }: { status: TicketService.Status }) {
	boardStore.update((state) => {
		const selected = cloneDeep(state.selected);
		if (selected && selected.id === status.board_id) {
			selected.statuses.push(status);
		}

		state.selected = selected;

		return state;
	});
}

function addTicket({ board_id, ticket }: { board_id: number; ticket: TicketService.Ticket }) {
	boardStore.update((state) => {
		const selected = cloneDeep(state.selected);
		if (selected && selected.id === board_id) {
			const status = selected.statuses.find((s) => s.id === ticket.status_id);
			if (status) {
				status.tickets.push(ticket);
			}
		}

		state.selected = selected;

		return state;
	});
}

export const BoardStore = {
	...boardStore,
	defaultState,
	selectBoard,
	addBoard,
	addStatus,
	addTicket
};
