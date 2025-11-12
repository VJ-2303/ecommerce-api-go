// Leaderboard Page

document.addEventListener('DOMContentLoaded', () => {
    setupMobileNav();
    loadLeaderboard();
});

async function loadLeaderboard() {
    const podium = document.getElementById('podium');
    const leaderboardBody = document.getElementById('leaderboardBody');

    if (!podium || !leaderboardBody) return;

    try {
        const data = await statsApi.getLeaderboard();

        if (data.leaderboard && data.leaderboard.length > 0) {
            displayPodium(data.leaderboard.slice(0, 3));
            displayLeaderboardTable(data.leaderboard);
        } else {
            podium.innerHTML = '<div class="loading">No leaderboard data available</div>';
            leaderboardBody.innerHTML = '<tr><td colspan="3" class="loading">No data available</td></tr>';
        }
    } catch (error) {
        console.error('Error loading leaderboard:', error);
        podium.innerHTML = '<div class="loading">Failed to load leaderboard</div>';
        leaderboardBody.innerHTML = '<tr><td colspan="3" class="loading">Failed to load data</td></tr>';
    }
}

function displayPodium(topThree) {
    const podium = document.getElementById('podium');
    
    if (topThree.length === 0) {
        podium.innerHTML = '<div class="loading">No top contributors yet</div>';
        return;
    }

    // Reorder for podium display: 2nd, 1st, 3rd
    const podiumOrder = [];
    if (topThree[1]) podiumOrder.push({ ...topThree[1], rank: 2 });
    if (topThree[0]) podiumOrder.push({ ...topThree[0], rank: 1 });
    if (topThree[2]) podiumOrder.push({ ...topThree[2], rank: 3 });

    podium.innerHTML = '';

    podiumOrder.forEach(entry => {
        const place = document.createElement('div');
        place.className = 'podium-place';
        
        const medals = { 1: '🥇', 2: '🥈', 3: '🥉' };
        
        place.innerHTML = `
            <div class="podium-rank">${medals[entry.rank]}</div>
            <div class="podium-name">${escapeHtml(entry.user_name)}</div>
            <div class="podium-count">${entry.report_count} reports</div>
        `;
        
        podium.appendChild(place);
    });
}

function displayLeaderboardTable(leaderboard) {
    const tbody = document.getElementById('leaderboardBody');
    
    tbody.innerHTML = '';

    leaderboard.forEach((entry, index) => {
        const row = document.createElement('tr');
        
        const rankClass = index < 3 ? 'rank-badge' : '';
        const rankEmoji = index === 0 ? '🥇' : index === 1 ? '🥈' : index === 2 ? '🥉' : '';
        
        row.innerHTML = `
            <td>
                ${rankClass ? `<span class="${rankClass}">${rankEmoji}</span>` : index + 1}
            </td>
            <td>
                <strong>${escapeHtml(entry.user_name)}</strong>
                <div style="font-size: 0.875rem; color: var(--text-secondary);">
                    ${escapeHtml(entry.phone_number)}
                </div>
            </td>
            <td><strong>${entry.report_count}</strong></td>
        `;
        
        tbody.appendChild(row);
    });
}

function setupMobileNav() {
    const navToggle = document.getElementById('navToggle');
    const navMenu = document.getElementById('navMenu');

    if (navToggle && navMenu) {
        navToggle.addEventListener('click', () => {
            navMenu.classList.toggle('active');
        });
    }
}
