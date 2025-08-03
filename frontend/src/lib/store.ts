import { writable } from 'svelte/store';

export type ViewState = 'setup' | 'dashboard' | 'loading';

// Create a writable store with an initial value of 'loading'
export const viewState = writable<ViewState>('setup');