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

    // Reorder for podium display: 2nd, 1st, 3rd (visual hierarchy)
    const podiumOrder = [];
    if (topThree[1]) podiumOrder.push({ ...topThree[1], rank: 2, displayRank: '🥈' });
    if (topThree[0]) podiumOrder.push({ ...topThree[0], rank: 1, displayRank: '🥇' });
    if (topThree[2]) podiumOrder.push({ ...topThree[2], rank: 3, displayRank: '🥉' });

    podium.innerHTML = '';

    podiumOrder.forEach(entry => {
        const place = document.createElement('div');
        place.className = 'podium-place';
        
        place.innerHTML = `
            <div class="podium-rank">${entry.displayRank}</div>
            <div class="podium-name">${escapeHtml(entry.user_name)}</div>
            <div class="podium-count">${entry.report_count} ${entry.report_count === 1 ? 'report' : 'reports'}</div>
        `;
        
        podium.appendChild(place);
    });
}

function displayLeaderboardTable(leaderboard) {
    const tbody = document.getElementById('leaderboardBody');
    
    tbody.innerHTML = '';

    leaderboard.forEach((entry, index) => {
        const row = document.createElement('tr');
        
        const medals = ['🥇', '🥈', '🥉'];
        const isTopThree = index < 3;
        
        row.innerHTML = `
            <td>
                ${isTopThree ? `<span class="rank-badge">${medals[index]}</span>` : `<span>${index + 1}</span>`}
            </td>
            <td>
                <strong>${escapeHtml(entry.user_name)}</strong>
                <div style="font-size: 0.8125rem; color: var(--text-secondary); margin-top: 2px;">
                    ${escapeHtml(entry.phone_number)}
                </div>
            </td>
            <td><strong>${entry.report_count}</strong> ${entry.report_count === 1 ? 'report' : 'reports'}</td>
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
