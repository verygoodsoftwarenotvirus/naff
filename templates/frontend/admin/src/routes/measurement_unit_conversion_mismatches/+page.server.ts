import type { PageServerLoad } from './$types';
import { getMeasurementUnitConversionMismatches } from '$lib/grpc/clients';

export const load: PageServerLoad = async ({ locals }) => {
  const token = locals.accessToken;
  if (!token) {
    return { items: [], error: 'Not authenticated' };
  }
  try {
    const res = (await getMeasurementUnitConversionMismatches(token, {})) as {
      mismatches?: Array<{
        ingredient?: { id?: string; name?: string };
        fromUnit?: { id?: string; name?: string };
        toUnit?: { id?: string; name?: string };
      }>;
    };
    return { items: res?.mismatches ?? [] };
  } catch (e) {
    return {
      items: [],
      error: e instanceof Error ? e.message : 'Failed to load conversion mismatches',
    };
  }
};
