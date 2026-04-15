<script lang="ts">
  import { Heading, Text, Card } from '@dinnerdonebetter/ui';

  let { data } = $props();
</script>

<Heading level={1}>Measurement Unit Conversion Mismatches</Heading>
<p class="subtitle">View and add conversion mismatches</p>

{#if data?.error}
  <p class="error">{data.error}</p>
{:else if data?.items && data.items.length > 0}
  <div class="table-wrap">
    <table class="data-table">
      <thead>
        <tr>
          <th>Ingredient</th>
          <th>From unit</th>
          <th>To unit</th>
        </tr>
      </thead>
      <tbody>
        {#each data.items as item, i ((item as Record<string, unknown>).id ?? i)}
          {@const row = item as Record<string, unknown>}
          {@const ing = row.ingredient as { id?: string; name?: string } | undefined}
          {@const from = row.fromUnit as { id?: string; name?: string } | undefined}
          {@const to = row.toUnit as { id?: string; name?: string } | undefined}
          <tr>
            <td>{ing?.name ?? '-'} <code class="id">{ing?.id ?? ''}</code></td>
            <td>{from?.name ?? '-'} <code class="id">{from?.id ?? ''}</code></td>
            <td>{to?.name ?? '-'} <code class="id">{to?.id ?? ''}</code></td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
{:else}
  <Card>
    <Text>No conversion mismatches found.</Text>
  </Card>
{/if}

<style>
  .subtitle {
    color: var(--color-text-muted);
    margin-bottom: var(--space-lg);
  }
  .error {
    color: var(--color-error, #b91c1c);
  }
  .table-wrap {
    overflow-x: auto;
  }
  .data-table {
    width: 100%;
    border-collapse: collapse;
  }
  .data-table th,
  .data-table td {
    padding: var(--space-sm) var(--space-md);
    text-align: left;
    border-bottom: 1px solid var(--color-border);
  }
  .id {
    font-size: 0.85em;
    color: var(--color-text-muted);
    margin-left: var(--space-xs);
  }
</style>
