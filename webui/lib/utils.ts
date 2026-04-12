import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function validateUrl(address: string) {
  try {
    new URL(address);
    return true;
  } catch (e) {
    return false;
  }
}

