<script lang="ts">
  import { Heading, Text, Link, Card } from '@dinnerdonebetter/ui';

  let { data } = $props();
  const item = data?.item as Record<string, unknown> | null;
</script>

{#if data?.error}
  <p class="error">{data.error}</p>
{:else if item}
  <Heading level={1}>Valid Ingredient</Heading>
  <Card>
    <dl class="detail-list">
      <dt>ID</dt>
      <dd><code>{item.id ?? '-'}</code></dd>
      <dt>Name</dt>
      <dd>{item.name ?? '-'}</dd>
    </dl>
  </Card>
  <p><Link href="/valid_ingredients">Back to valid ingredients</Link></p>
{:else}
  <Text>Ingredient not found.</Text>
{/if}

<style>
  .error {
    color: var(--color-error, #b91c1c);
  }
  .detail-list {
    display: grid;
    grid-template-columns: auto 1fr;
    gap: var(--space-xs) var(--space-lg);
  }
  .detail-list dt {
    color: var(--color-text-muted);
  }
</style>
