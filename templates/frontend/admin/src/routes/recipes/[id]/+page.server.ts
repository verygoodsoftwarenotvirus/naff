import type { PageServerLoad } from './$types';
import { getRecipe } from '$lib/grpc/clients';

export const load: PageServerLoad = async ({ locals, params }) => {
  const token = locals.accessToken;
  const id = params.id;
  if (!token) {
    return { recipe: null, error: 'Not authenticated' };
  }
  try {
    const res = (await getRecipe(token, { recipeId: id })) as { result?: Record<string, unknown> };
    return { recipe: res?.result ?? null };
  } catch (e) {
    return {
      recipe: null,
      error: e instanceof Error ? e.message : 'Failed to load recipe',
    };
  }
};
