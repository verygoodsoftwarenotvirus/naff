import { buildServerSideLogger } from '@dinnerdonebetter/logger';

const isDev = process.env.NODE_ENV !== 'production';

export const logger = buildServerSideLogger('admin', isDev);
