<script lang="ts">
  import { Search, Group, SunMoon, CircleUser } from '@lucide/svelte';

  let { searchTerm = $bindable(''), isGrouped = $bindable(false), username = '' } = $props();

  let html = document.documentElement;
  let currentTheme = $state(localStorage.getItem('quadboard-theme') || 'light');
  
  // Variable pour garder la barre ouverte si on tape dedans
  let isFocused = $state(false);

  $effect(() => {
    html.setAttribute('data-theme', currentTheme);
    
    const favicon = document.getElementById('favicon') as HTMLLinkElement | null;
    
    if (favicon) {
      favicon.href = currentTheme === 'dark' ? '/assets/img/quadboard-bw.svg' : '/assets/img/quadboard-color.svg';
    }
  });

  function toggleTheme() {
    currentTheme = currentTheme === 'light' ? 'dark' : 'light';
    localStorage.setItem('quadboard-theme', currentTheme);
  }
</script>

<header class="header">
  <div class="brand">
    <div class="logo-container">
      <img src="/assets/img/quadboard-color.svg" alt="QuadBoard" class="logo logo-light">
      <img src="/assets/img/quadboard-bw.svg" alt="QuadBoard" class="logo logo-dark">
    </div>
    <h1>QuadBoard</h1>
  </div>
  
  <div class="header-actions">
    <!-- On ajoute la classe 'expanded' si focus ou s'il y a du texte -->
    <div class="search-box" class:expanded={isFocused || searchTerm.length > 0}>
      <input 
        type="text" 
        class="search-input"
        bind:value={searchTerm} 
        onfocus={() => isFocused = true}
        onblur={() => isFocused = false}
        placeholder="Search services..." 
        autocomplete="off"
      >
      <span class="search-icon">
        <Search size={18} />
      </span>
    </div>
    
    <button class="btn-icon" class:active={isGrouped} onclick={() => isGrouped = !isGrouped} aria-label="Toggle grouping">
      <Group size={20} />
    </button>

    <button class="btn-icon" onclick={toggleTheme} aria-label="Toggle theme">
      <SunMoon size={20} />
    </button>

    {#if username}
      <div class="user-profile">
        <CircleUser size={20} />
        <span>{username}</span>
      </div>
    {/if}
  </div>
</header>

<style>
  .header {
    display: flex; justify-content: space-between; align-items: center;
    padding: 1rem 2rem; background: var(--header-bg);
    backdrop-filter: blur(10px); border-bottom: 1px solid var(--border-color);
    position: sticky; top: 0; z-index: 10;
  }
  .brand { display: flex; align-items: center; gap: 0.75rem; }
  .logo-container {
    width: 36px; height: 36px; background-color: #ffffff;
    border-radius: 8px; display: flex; align-items: center; justify-content: center;
    padding: 4px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); overflow: hidden;
  }
  .logo { width: 100%; height: 100%; object-fit: contain; display: block; }
  .logo-dark { display: none; }
  :global([data-theme="dark"]) .logo-light { display: none; }
  :global([data-theme="dark"]) .logo-dark { display: block; }
  
  h1 { font-size: 1.25rem; font-weight: 700; letter-spacing: -0.025em; }
  .header-actions { display: flex; align-items: center; gap: 1rem; }
  
  .search-box { 
    position: relative; 
    display: flex; 
    align-items: center;
  }
  .search-input {
    width: 38px; 
    height: 38px; 
    padding-left: 12px;
    padding-right: 30px; 
    font-size: 14px;
    color: var(--text-color); 
    background: var(--card-bg);
    border: 1px solid var(--border-color); 
    border-radius: 30px;
    transition: all 0.4s ease; 
    outline: none; 
    cursor: pointer;
  }
  .search-input::placeholder { color: transparent; transition: color 0.3s ease 0.1s; }
  
  /* La barre s'étend au survol ou si la classe 'expanded' est active */
  .search-box:hover .search-input, 
  .search-box.expanded .search-input {
    width: 220px; 
    border-color: var(--primary-color);
    box-shadow: 0 0 10px rgba(37, 99, 235, 0.2); 
    cursor: text;
  }
  
  .search-box:hover .search-input::placeholder, 
  .search-box.expanded .search-input::placeholder {
    color: var(--text-muted);
  }
  
  .search-icon {
    position: absolute; 
    right: 12px; 
    top: 50%; /* Centrage vertical */
    transform: translateY(-50%); /* Centrage vertical */
    color: var(--text-muted); 
    pointer-events: none; 
    transition: color 0.3s ease;
  }
  
  .search-box:hover .search-icon, 
  .search-box.expanded .search-icon { 
    color: var(--primary-color); 
  }
  
  .btn-icon {
    background: none; border: 1px solid var(--border-color); border-radius: 8px;
    padding: 0.5rem; cursor: pointer; color: var(--text-color);
    display: flex; align-items: center; justify-content: center; transition: background-color 0.2s;
  }
  .btn-icon:hover { background-color: var(--input-bg); }
  .btn-icon.active { background-color: var(--primary-color); border-color: var(--primary-color); color: white; }
  
  .user-profile { display: flex; align-items: center; gap: 0.5rem; font-size: 0.875rem; font-weight: 500; color: var(--text-muted); }
</style>