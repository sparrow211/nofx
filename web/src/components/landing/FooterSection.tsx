import { t, Language } from '../../i18n/translations'

interface FooterSectionProps {
  language: Language
}

export default function FooterSection({ language }: FooterSectionProps) {
  return (
    <footer style={{ borderTop: '1px solid var(--panel-border)', background: 'var(--brand-dark-gray)' }}>
      <div className='max-w-[1200px] mx-auto px-6 py-10'>
        {/* Brand */}
        <div className='flex items-center gap-3 mb-8'>
          <img src='/icons/aibtc.svg' alt='NF Logo' className='w-8 h-8' />
          <div>
            <div className='text-lg font-bold' style={{ color: '#EAECEF' }}>
              NF
            </div>
            <div className='text-xs' style={{ color: '#848E9C' }}>
              {t('futureStandardAI', language)}
            </div>
          </div>
        </div>

        {/* Multi-link columns */}

        {/* Bottom note (kept subtle) */}
        <div
          className='pt-6 mt-8 text-center text-xs'
          style={{ color: 'var(--text-tertiary)', borderTop: '1px solid var(--panel-border)' }}
        >
          <p>{t('footerTitle', language)}</p>
          <p className='mt-1'>{t('footerWarning', language)}</p>
        </div>
      </div>
    </footer>
  )
}
