import { create } from "zustand";

export interface AppState {}

const useAppStore = create<AppState>((set) => ({}));
