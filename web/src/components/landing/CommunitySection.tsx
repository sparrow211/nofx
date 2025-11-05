import { motion } from 'framer-motion'
import AnimatedSection from './AnimatedSection'

interface CardProps {
  quote: string;
  authorName: string;
  handle: string;
  avatarUrl: string;
  tweetUrl: string;
  delay: number;
}

function TestimonialCard({ quote, authorName, delay }: CardProps) {
  return (
    <motion.div
      className='p-6 rounded-xl'
      style={{ background: 'var(--brand-dark-gray)', border: '1px solid rgba(240, 185, 11, 0.1)' }}
      initial={{ opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true }}
      transition={{ delay }}
      whileHover={{ scale: 1.05 }}
    >
      <p className='text-lg mb-4' style={{ color: 'var(--brand-light-gray)' }}>
        "{quote}"
      </p>
      <div className='flex items-center gap-2'>
        <div className='w-8 h-8 rounded-full' style={{ background: 'var(--binance-yellow)' }} />
        <span className='text-sm font-semibold' style={{ color: 'var(--text-secondary)' }}>
          {authorName}
        </span>
      </div>
    </motion.div>
  )
}

export default function CommunitySection() {
  const staggerContainer = { animate: { transition: { staggerChildren: 0.1 } } }

  // 推特内容整合（保持原三列布局，超出自动换行）
  const items: CardProps[] = [
    {
      quote:
        '昨夜 AIBTC 自动交易忙，AI 似灵狐，于多空间狡黠游走。无需深谙交易理，财富悄然涨，一晚揽下 6% 盈光，开启躺赚新章！',
      authorName: 'Panda',
      handle: '@Panda',
      avatarUrl:
        'https://q1.qlogo.cn/g?b=qq&nk=414167924&s=640',
      tweetUrl:
        '#',
      delay: 0,
    },
    {
      quote:
        ' AI 自动交易，太有意思了，就看 AI 在那一会开空一会开多，一顿操作，虽然看不懂为什么，但是一晚上帮我赚了 6% 收益',
      authorName: 'ddd',
      handle: '@ddd',
      avatarUrl:
        'https://q1.qlogo.cn/g?b=qq&nk=414167955&s=640',
      tweetUrl: '#',
      delay: 0.1,
    },
    {
      quote:
        'Dive into an all - night AIBTC auto - trading extravaganza! The AI transforms into a trading sprite, gracefully flitting between long and short markets with each move resembling a mysterious dance. Even if you don\'t grasp the trading logic, you can just sit back and reap the rewards. In a single night, it delivers a solid 6% return, making wealth growth a breeze!',
      authorName: 'Jay',
      handle: '@jay',
      avatarUrl:
        'https://q1.qlogo.cn/g?b=qq&nk=9851356&s=640',
      tweetUrl: '#',
      delay: 0.15,
    },
  ]

  return (
    <AnimatedSection>
      <div className='max-w-7xl mx-auto'>
        <motion.div
          className='grid md:grid-cols-3 gap-6'
          variants={staggerContainer}
          initial='initial'
          whileInView='animate'
          viewport={{ once: true }}
        >
          {items.map((item, idx) => (
            <TestimonialCard key={idx} {...item} />
          ))}
        </motion.div>
      </div>
    </AnimatedSection>
  )
}
