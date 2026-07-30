<script lang="ts">
	import { Lock, Box, ExternalLink } from "@lucide/svelte";

	let { app, groupName, onGroupFilter } = $props();

	const icons: Record<string, any> = { lock: Lock, box: Box };

	function handleLogoError(e: Event) {
		const target = e.target as HTMLImageElement;
		target.style.display = "none";
		if (target.nextElementSibling) {
			(target.nextElementSibling as HTMLElement).style.display = "inline-block";
		}
	}

	const iconKey = $derived((app.Icon as string) || "box");
	const DynIcon = $derived(icons[iconKey] || Box);

	function handleBadgeClick(e: MouseEvent) {
		e.preventDefault();
		e.stopPropagation();
		onGroupFilter(groupName);
	}
</script>

{#if !app.Authorized}
	<div class="card disabled">
		<div class="card-content">
			<div class="card-header">
				<div class="card-icon-wrapper">
					<Lock class="card-icon" size={20} />
				</div>
				<h3 class="card-title">{app.Name}</h3>
			</div>
		</div>
		<div class="card-footer">
			<button class="card-group-badge" onclick={handleBadgeClick}
				>{groupName}</button
			>
			<span class="card-no-url">Restricted</span>
		</div>
	</div>
{:else}
	<a
		href={app.URL || "#"}
		target="_blank"
		rel="noopener noreferrer"
		class="card"
		class:disabled={!app.URL}
	>
		<div class="card-content">
			<div class="card-header">
				<div class="card-icon-wrapper">
					{#if app.Logo}
						<img
							src={app.Logo}
							alt={app.Name}
							class="card-logo"
							onerror={handleLogoError}
						/>
						<DynIcon class="card-icon" style="display: none;" size={20} />
					{:else}
						<DynIcon class="card-icon" size={20} />
					{/if}
				</div>
				<h3 class="card-title">{app.Name}</h3>
			</div>
			{#if app.Description}
				<p class="card-desc">{app.Description}</p>
			{/if}
		</div>
		<div class="card-footer">
			<button class="card-group-badge" onclick={handleBadgeClick}
				>{groupName}</button
			>
			{#if app.URL}
				<span class="card-open">Open <ExternalLink size={12} /></span>
			{:else}
				<span class="card-no-url">No URL</span>
			{/if}
		</div>
	</a>
{/if}

<style>
	.card {
		background-color: var(--card-bg);
		border: 1px solid var(--border-color);
		border-radius: 12px;
		padding: 1.25rem;
		text-decoration: none;
		color: inherit;
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		box-shadow: var(--shadow);
		transition:
			transform 0.2s ease,
			box-shadow 0.2s ease,
			border-color 0.2s,
			opacity 0.2s;
		overflow: hidden;
	}
	.card:hover {
		transform: translateY(-4px);
		border-color: var(--primary-color);
		box-shadow:
			0 10px 15px -3px rgba(0, 0, 0, 0.1),
			0 4px 6px -2px rgba(0, 0, 0, 0.05);
	}
	.card.disabled {
		opacity: 0.6;
		pointer-events: none;
	}

	.card-header {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		margin-bottom: 0.75rem;
	}
	.card-icon-wrapper {
		width: 36px;
		height: 36px;
		background-color: var(--input-bg);
		border: 1px solid var(--border-color);
		border-radius: 8px;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
	}
	:global(.card-icon) {
		width: 20px;
		height: 20px;
		color: var(--primary-color);
	}
	.card-logo {
		width: 20px;
		height: 20px;
		object-fit: contain;
		max-width: 100%;
	}
	.card-title {
		font-size: 1rem;
		font-weight: 600;
		color: var(--text-color);
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	.card-desc {
		color: var(--text-muted);
		font-size: 0.8rem;
		line-height: 1.4;
		display: -webkit-box;
		line-clamp: 2;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}

	.card-footer {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-top: 1rem;
		font-size: 0.75rem;
		color: var(--text-muted);
	}
	.card-group-badge {
		display: inline-flex;
		align-items: center;
		padding: 0.25rem 0.6rem;
		font-size: 0.7rem;
		font-weight: 600;
		border-radius: 9999px;
		background-color: rgba(16, 185, 129, 0.15);
		color: var(--secondary-color);
		border: 1px solid transparent;
		cursor: pointer;
		transition: all 0.2s ease;
	}
	.card-group-badge:hover {
		background-color: var(--secondary-color);
		color: white;
	}

	.card-open {
		display: flex;
		align-items: center;
		gap: 0.25rem;
		font-weight: 600;
		color: var(--primary-color);
	}
	.card-no-url {
		font-style: italic;
		color: var(--text-muted);
	}
</style>
