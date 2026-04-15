import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { env } from '$env/dynamic/private';

export const GET: RequestHandler = async () => {
  const teamId = env.AASA_TEAM_ID ?? env.DINNER_DONE_BETTER_AASA_TEAM_ID ?? '';
  const bundleId = env.AASA_BUNDLE_ID ?? env.DINNER_DONE_BETTER_AASA_BUNDLE_ID ?? '';
  const appId = teamId && bundleId ? `${teamId}.${bundleId}` : '';

  const aasa = {
    applinks: {
      apps: [] as string[],
      details: [
        {
          appID: appId,
          paths: ['/accept_invitation', '/accept_invitation/*', '/meal_plans', '/meal_plans/*'],
        },
      ],
    },
    webcredentials: {
      apps: appId ? [appId] : [],
    },
  };

  return json(aasa, {
    headers: {
      'Content-Type': 'application/json',
    },
  });
};
