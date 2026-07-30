<script lang="ts">
	import { onMount } from "svelte";
	import Header from "./lib/Header.svelte";
	import Card from "./lib/Card.svelte";
	import { SearchX } from "@lucide/svelte";

	let resources = $state([] as any[]);
	let searchTerm = $state("");
	let isGrouped = $state(false);
	let username = $state(""); // <-- Ajouté
	let activeGroupFilter = $state<string | null>(null);

	onMount(async () => {
		try {
			// 1. Récupération des infos utilisateur et config
			const meRes = await fetch("/api/v1/me");
			if (meRes.ok) {
				const meData = await meRes.json();
				username = meData.username || "";
				isGrouped = meData.groupByDefault;
			}

			// 2. Récupération du catalogue
			const res = await fetch("/api/v1/catalog");
			resources = await res.json();
		} catch (e) {
			console.error(e);
		}
	});

	let filtered = $derived(
		resources.filter((app) => {
			const text = (app.Name + " " + (app.Description || "")).toLowerCase();
			return text.includes(searchTerm.toLowerCase());
		}),
	);

	let grouped = $derived(
		filtered.reduce((acc: Record<string, any[]>, app) => {
			const name = app.DisplayName || app.Group || "Default";
			if (!acc[name]) acc[name] = [];
			acc[name].push(app);
			return acc;
		}, {}),
	);

	let displayResources = $derived(
		activeGroupFilter
			? filtered.filter((app) => (app.Group || "Default") === activeGroupFilter)
			: filtered,
	);
</script>

<!-- On passe username en prop au Header -->
<Header bind:searchTerm bind:isGrouped {username} />

<main class="container">
	{#if resources.length === 0}
		<div class="card-grid">
			{#each Array(8) as _, i}
				<div class="card skeleton">
					<div class="card-content">
						<div class="card-header">
							<div class="card-icon-wrapper skeleton-bg">
								<div class="skeleton-avatar"></div>
							</div>
							<div class="skeleton-line title"></div>
						</div>
						<div class="skeleton-line desc"></div>
						<div class="skeleton-line desc short"></div>
					</div>
					<div class="card-footer"><div class="skeleton-line link"></div></div>
				</div>
			{/each}
		</div>
	{:else if filtered.length === 0}
		<div class="empty-state">
			<SearchX size={48} />
			<p>Aucune application trouvée.</p>
		</div>
	{:else}
		{#if activeGroupFilter}
			<div class="filter-banner">
				<span>Filtre actif : <strong>{activeGroupFilter}</strong></span>
				<button onclick={() => (activeGroupFilter = null)}
					>Effacer le filtre</button
				>
			</div>
		{/if}

    {#if isGrouped && !activeGroupFilter}
      {#each Object.keys(grouped).sort() as groupName}
        <section class="resource-group">
          <h2 class="group-title">{groupName}</h2>
          <div class="card-grid">
            {#each grouped[groupName] as app}
              <!-- Ajout de : string -->
              <Card {app} {groupName} onGroupFilter={(g: string) => activeGroupFilter = g} />
            {/each}
          </div>
        </section>
      {/each}
    {:else}
      <div class="card-grid">
        {#each displayResources as app}
          <!-- Ajout de : string -->
          <Card {app} groupName={app.DisplayName || app.Group || 'Default'} onGroupFilter={(g: string) => activeGroupFilter = g} />
        {/each}
      </div>
    {/if}
	{/if}
</main>

<style>
	.container {
		max-width: 1280px;
		margin: 2rem auto;
		padding: 0 2rem;
	}
	.resource-group {
		margin-bottom: 2.5rem;
	}
	.group-title {
		font-size: 1.1rem;
		font-weight: 600;
		color: var(--text-muted);
		margin-bottom: 1rem;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		border-bottom: 1px solid var(--border-color);
		padding-bottom: 0.5rem;
	}
	.card-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
		gap: 1.5rem;
	}

	.filter-banner {
		display: flex;
		justify-content: space-between;
		align-items: center;
		background-color: var(--card-bg);
		border: 1px solid var(--border-color);
		border-radius: 8px;
		padding: 0.75rem 1rem;
		margin-bottom: 1.5rem;
		font-size: 0.9rem;
		color: var(--text-muted);
	}
	.filter-banner button {
		background: none;
		border: 1px solid var(--border-color);
		border-radius: 6px;
		padding: 0.3rem 0.6rem;
		cursor: pointer;
		color: var(--text-color);
		font-size: 0.8rem;
	}

	.card.skeleton {
		pointer-events: none;
	}
	.card-icon-wrapper.skeleton-bg {
		background-color: var(--input-bg);
	}
	.skeleton-avatar {
		width: 20px;
		height: 20px;
		border-radius: 4px;
		background-color: var(--border-color);
		animation: pulse 1.5s infinite ease-in-out;
	}
	.skeleton-line {
		height: 12px;
		background-color: var(--border-color);
		border-radius: 4px;
		animation: pulse 1.5s infinite ease-in-out;
	}
	.skeleton-line.title {
		width: 120px;
	}
	.skeleton-line.desc {
		margin-top: 0.75rem;
		width: 100%;
	}
	.skeleton-line.desc.short {
		margin-top: 0.5rem;
		width: 60%;
	}
	.skeleton-line.link {
		width: 60px;
		height: 16px;
	}
	@keyframes pulse {
		0% {
			opacity: 0.6;
		}
		50% {
			opacity: 0.3;
		}
		100% {
			opacity: 0.6;
		}
	}

	.empty-state {
		grid-column: 1 / -1;
		text-align: center;
		color: var(--text-muted);
		padding: 4rem 0;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
	}
</style>
