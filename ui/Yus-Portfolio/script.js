// YUS – Yelloh Bus | Vanilla JavaScript

document.addEventListener('DOMContentLoaded', () => {
    // Initialize all modules
    initScrollProgress();
    initScrollAnimations();
    initRouteProgress();
    initPassengerAnimation();
    initFinalBusAnimation();
});

// Scroll Progress Bar
function initScrollProgress() {
    const progressBar = document.getElementById('progressBar');
    const scrollBus = document.getElementById('scrollBus');
    const scrollNodes = document.querySelectorAll('.scroll-node');
    
    if (!progressBar) return;
    
    function updateProgress() {
        const scrollTop = window.pageYOffset;
        const docHeight = document.documentElement.scrollHeight - window.innerHeight;
        const scrollPercent = (scrollTop / docHeight) * 100;
        
        // Update progress bar
        progressBar.style.width = scrollPercent + '%';
        
        // Update bus position
        if (scrollBus) {
            const maxMove = window.innerWidth - 150;
            const busMove = (scrollPercent / 100) * maxMove;
            scrollBus.style.transform = `translateX(${busMove}px)`;
        }
        
        // Update nodes
        scrollNodes.forEach((node, index) => {
            const nodePercent = (index / (scrollNodes.length - 1)) * 100;
            if (scrollPercent >= nodePercent - 5) {
                node.classList.add('active');
            } else {
                node.classList.remove('active');
            }
        });
    }
    
    window.addEventListener('scroll', updateProgress, { passive: true });
    updateProgress();
}

// Scroll-triggered Animations
function initScrollAnimations() {
    const animatedElements = document.querySelectorAll('[data-animate]');
    
    const observerOptions = {
        threshold: 0.1,
        rootMargin: '0px 0px -50px 0px'
    };
    
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const delay = entry.target.dataset.delay || 0;
                setTimeout(() => {
                    entry.target.classList.add('visible');
                }, parseInt(delay));
            }
        });
    }, observerOptions);
    
    animatedElements.forEach(el => observer.observe(el));
}

// Route Section Progress Animation
function initRouteProgress() {
    const routeSection = document.getElementById('route');
    const routeProgress = document.getElementById('routeProgress');
    const routeNodes = document.querySelectorAll('.route-node');
    
    if (!routeSection || !routeProgress) return;
    
    function updateRouteProgress() {
        const rect = routeSection.getBoundingClientRect();
        const sectionHeight = routeSection.offsetHeight;
        const viewportHeight = window.innerHeight;
        
        // Calculate how much of the section is visible
        const visibleStart = Math.max(0, -rect.top);
        const visibleEnd = Math.min(sectionHeight, viewportHeight - rect.top);
        const visibleAmount = visibleEnd - visibleStart;
        
        // Calculate progress percentage
        let progress = 0;
        if (rect.top < viewportHeight && rect.bottom > 0) {
            progress = Math.min(100, Math.max(0, ((viewportHeight - rect.top) / (sectionHeight + viewportHeight)) * 100));
        }
        
        routeProgress.style.height = progress + '%';
        
        // Activate nodes based on progress
        routeNodes.forEach((node, index) => {
            const nodePercent = ((index + 1) / routeNodes.length) * 100;
            if (progress >= nodePercent * 0.8) {
                node.classList.add('active');
            } else if (index > 0) {
                node.classList.remove('active');
            }
        });
    }
    
    window.addEventListener('scroll', updateRouteProgress, { passive: true });
    updateRouteProgress();
}

// Passenger Section Animation
function initPassengerAnimation() {
    const passengerStops = document.getElementById('passengerStops');
    const etaTime = document.getElementById('etaTime');
    
    if (!passengerStops) return;
    
    const stops = passengerStops.querySelectorAll('.passenger-stop');
    const etas = ['2 min', '7 min', '12 min', '17 min', '22 min'];
    let currentStop = 0;
    let isInView = false;
    
    // Intersection Observer for passenger section
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            isInView = entry.isIntersecting;
        });
    }, { threshold: 0.4 });
    
    const passengerSection = document.getElementById('passenger');
    if (passengerSection) {
        observer.observe(passengerSection);
    }
    
    function updateCurrentStop() {
        if (!isInView) return;
        
        // Remove current from all
        stops.forEach((stop, index) => {
            stop.classList.remove('current', 'passed');
            const node = stop.querySelector('.passenger-node');
            node.classList.remove('current');
            
            // Find and remove NOW badge
            const nowBadge = stop.querySelector('.now-badge');
            if (nowBadge) nowBadge.remove();
            
            // Update status
            if (index < currentStop) {
                stop.classList.add('passed');
                const infoP = stop.querySelector('.passenger-stop-info p');
                if (infoP) infoP.textContent = 'Passed';
            } else if (index === currentStop) {
                stop.classList.add('current');
                node.classList.add('current');
                
                // Add NOW badge
                const badge = document.createElement('span');
                badge.className = 'now-badge';
                badge.textContent = 'NOW';
                stop.appendChild(badge);
                
                // Update mini bus in node
                node.innerHTML = `
                    <svg viewBox="0 0 120 60" fill="none" xmlns="http://www.w3.org/2000/svg" class="mini-bus">
                        <rect x="5" y="10" width="100" height="35" rx="8" fill="#FBBF24"/>
                        <circle cx="28" cy="48" r="8" fill="#1e2128"/>
                        <circle cx="82" cy="48" r="8" fill="#1e2128"/>
                    </svg>
                    <div class="pulse-ring"></div>
                `;
            } else {
                // Reset to pin icon
                node.innerHTML = `
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"/>
                        <circle cx="12" cy="10" r="3"/>
                    </svg>
                `;
                const infoP = stop.querySelector('.passenger-stop-info p');
                if (infoP) infoP.textContent = etas[index];
            }
        });
        
        // Update ETA display
        if (etaTime) {
            etaTime.textContent = etas[currentStop];
            etaTime.style.animation = 'none';
            etaTime.offsetHeight; // Trigger reflow
            etaTime.style.animation = 'fadeIn 0.3s ease';
        }
        
        // Update progress line
        const progressLine = document.querySelector('.passenger-line-progress');
        if (progressLine) {
            const progressPercent = (currentStop / (stops.length - 1)) * 100;
            progressLine.style.height = progressPercent + '%';
        }
        
        currentStop = (currentStop + 1) % stops.length;
    }
    
    // Run animation every 2 seconds
    setInterval(updateCurrentStop, 2000);
    
    // Add CSS for fade animation
    const style = document.createElement('style');
    style.textContent = `
        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(10px); }
            to { opacity: 1; transform: translateY(0); }
        }
    `;
    document.head.appendChild(style);
}

// Final Bus Animation
function initFinalBusAnimation() {
    const finalBus = document.getElementById('finalBus');
    const finalSection = document.getElementById('final');
    
    if (!finalBus || !finalSection) return;
    
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                setTimeout(() => {
                    finalBus.classList.add('visible');
                }, 300);
            }
        });
    }, { threshold: 0.3 });
    
    observer.observe(finalSection);
}

// Smooth scroll for anchor links (if any are added later)
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function(e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href'));
        if (target) {
            target.scrollIntoView({
                behavior: 'smooth',
                block: 'start'
            });
        }
    });
});

// Parallax effect for ambient glows
function initParallax() {
    const glows = document.querySelectorAll('.ambient-glow');
    
    window.addEventListener('scroll', () => {
        const scrollY = window.pageYOffset;
        glows.forEach(glow => {
            const speed = 0.3;
            glow.style.transform = `translate(-50%, calc(-50% + ${scrollY * speed}px))`;
        });
    }, { passive: true });
}

// Initialize parallax on load
initParallax();

// Add hover effects for interactive elements
document.querySelectorAll('.problem-card, .stat-card, .app-card').forEach(card => {
    card.addEventListener('mouseenter', function() {
        this.style.transform = 'translateY(-5px)';
    });
    
    card.addEventListener('mouseleave', function() {
        this.style.transform = 'translateY(0)';
    });
});

// Console welcome message
console.log('%c🚌 YUS – Yelloh Bus Ecosystem', 'color: #FBBF24; font-size: 20px; font-weight: bold;');
console.log('%cA real-time college bus tracking ecosystem', 'color: #888; font-size: 14px;');
