<script lang="ts">
  import { Heading, Text, Link, Card } from '@dinnerdonebetter/ui';

  let { data } = $props();
  const recipe = data?.recipe as Record<string, unknown> | null;
</script>

{#if data?.error}
  <p class="error">{data.error}</p>
{:else if recipe}
  <Heading level={1}>Recipe</Heading>
  <Card>
    <dl class="detail-list">
      <dt>ID</dt>
      <dd><code>{recipe.id ?? '-'}</code></dd>
      <dt>Name</dt>
      <dd>{recipe.name ?? '-'}</dd>
    </dl>
  </Card>
  <p><Link href="/recipes">Back to recipes</Link></p>
{:else}
  <Text>Recipe not found.</Text>
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
