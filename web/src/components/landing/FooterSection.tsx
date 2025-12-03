import { t, Language } from '../../i18n/translations'

interface FooterSectionProps {
  language: Language
}

export default function FooterSection({ language }: FooterSectionProps) {
  return (
    <footer
      style={{
        borderTop: '1px solid var(--panel-border)',
        background: 'var(--brand-dark-gray)',
      }}
    >
      <div className="max-w-[1200px] mx-auto px-6 py-10">
        {/* Brand */}
        <div className="flex items-center gap-3 mb-8">
          <img src="/icons/aibtc.svg" alt="AIBTC Logo" className="w-8 h-8" />
          <div>
            <div className="text-lg font-bold" style={{ color: '#EAECEF' }}>
              AIBTC
            </div>
            <div className="text-xs" style={{ color: '#848E9C' }}>
              {t('futureStandardAI', language)}
            </div>
          </div>
        </div>

        {/* Multi-link columns */}
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-3 gap-8">
          <div>
            <h3
              className="text-sm font-semibold mb-3"
              style={{ color: '#EAECEF' }}
            >
              {t('links', language)}
            </h3>
            <ul className="space-y-2 text-sm" style={{ color: '#848E9C' }}>
              <li>
                <a
                  className="hover:text-[#F0B90B]"
                  href="#"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  xxx
                </a>
              </li>
            </ul>
          </div>

          <div>
            <h3
              className="text-sm font-semibold mb-3"
              style={{ color: '#EAECEF' }}
            >
              {t('resources', language)}
            </h3>
            <ul className="space-y-2 text-sm" style={{ color: '#848E9C' }}>
              {/* Documentation and issue trackers removed (external GitHub links) */}
            </ul>
          </div>

          <div>
            <h3
              className="text-sm font-semibold mb-3"
              style={{ color: '#EAECEF' }}
            >
              {t('supporters', language)}
            </h3>
            <ul className="space-y-2 text-sm" style={{ color: '#848E9C' }}>
              <li>
                <a
                  className="hover:text-[#F0B90B]"
                  href="https://s.qiniu.com/bya2Ub"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  1000万大模型免费领
                </a>
              </li>
              <li>
                <a
                  className="hover:text-[#F0B90B]"
                  href="https://s.qiniu.com/bya2Ub"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  1000万大模型免费领
                </a>
              </li>
              
            </ul>
          </div>
        </div>

        {/* Bottom note (kept subtle) */}
        <div
          className="pt-6 mt-8 text-center text-xs"
          style={{
            color: 'var(--text-tertiary)',
            borderTop: '1px solid var(--panel-border)',
          }}
        >
          <p>{t('footerTitle', language)}</p>
          <p className="mt-1">{t('footerWarning', language)}</p>
        </div>
      </div>
    </footer>
  )
}
