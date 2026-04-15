import type { PageServerLoad } from './$types';
import { getValidIngredient } from '$lib/grpc/clients';

export const load: PageServerLoad = async ({ locals, params }) => {
  const token = locals.accessToken;
  const id = params.id;
  if (!token) {
    return { item: null, error: 'Not authenticated' };
  }
  try {
    const res = (await getValidIngredient(token, { validIngredientId: id })) as {
      result?: Record<string, unknown>;
    };
    return { item: res?.result ?? null };
  } catch (e) {
    return {
      item: null,
      error: e instanceof Error ? e.message : 'Failed to load ingredient',
    };
  }
};
